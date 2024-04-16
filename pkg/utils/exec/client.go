package exec

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// Ensure all clients implement Client interface.
var _ Client = LocalClient{}
var _ Client = RemoteClient{}

type Client interface {
	Run(Command) error
	RunCtx(context.Context, Command) error
	Close() error
}

// commonClient provides functions that are common for all clients.
type commonClient struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (c *commonClient) SetStdin(stdin io.Reader) {
	c.stdin = stdin
}

func (c *commonClient) SetStdout(stdout io.Writer) {
	c.stdout = stdout
}

func (c *commonClient) SetStderr(stderr io.Writer) {
	c.stderr = stderr
}

func (c *commonClient) SetCombinedStdout(writer io.Writer) {
	c.stdout = writer
	c.stderr = writer
}

func (c commonClient) Close() error {
	return nil
}

type LocalClient struct {
	commonClient
}

// NewLocalClient initializes a client for running local commands.
func NewLocalClient() LocalClient {
	return LocalClient{}
}

// Run runs command locally.
func (c LocalClient) Run(command Command) error {
	return c.RunCtx(context.Background(), command)
}

// RunCtx runs command locally.
func (c LocalClient) RunCtx(ctx context.Context, command Command) error {
	cmd := exec.CommandContext(ctx, command.command, command.args...)
	cmd.Stdin = c.stdin
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr

	if len(command.envs) > 0 {
		env := make([]string, 0, len(command.envs))
		for k, v := range command.envs {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}

		cmd.Env = env
	}

	if command.workingDir != "" {
		cmd.Dir = command.workingDir
	}

	return cmd.Run()
}

type RemoteClient struct {
	commonClient

	user           string
	host           string
	port           string
	privateKeyPath string
	publicKeyPath  string
	initialized    bool
	sudo           bool

	client *ssh.Client
	mux    *sync.Mutex
}

// NewSSHClient initializes a new remote SSH client.
func NewSSHClient(user string, host string) RemoteClient {
	return RemoteClient{
		user: user,
		host: host,
		port: "22",
		mux:  &sync.Mutex{},
	}
}

// WithPort sets the host port to given value. By default, port 22 is used.
func (c RemoteClient) WithPort(port uint16) RemoteClient {
	// Immutable once client is initialized.
	if !c.isInitialized() {
		c.port = fmt.Sprint(port)
	}
	return c
}

// WithPrivateKeyFile sets the path to the private key file that is used
// for authentication..
func (c RemoteClient) WithPrivateKeyFile(privateKeyPath string) RemoteClient {
	// Immutable once client is initialized.
	if !c.isInitialized() {
		c.privateKeyPath = privateKeyPath
	}
	return c
}

// WithPublicKeyFile sets the path to the public key file that is used
// for host verification. If not set, known hosts are ignored.
func (c RemoteClient) WithPublicKeyFile(publicKeyPath string) RemoteClient {
	// Immutable once client is initialized.
	if !c.isInitialized() {
		c.publicKeyPath = publicKeyPath
	}
	return c
}

func (c RemoteClient) WithSuperUser(sudo bool) RemoteClient {
	// Immutable once client is initialized.
	if !c.isInitialized() {
		c.sudo = sudo
	}
	return c
}

// Endpoint returns SSH endpoint in format "user@host:port".
func (c RemoteClient) Endpoint() string {
	return fmt.Sprintf("%s@%s:%s", c.user, c.host, c.port)
}

// Close closes potentially initialized the SSH client.
func (c RemoteClient) Close() error {
	if !c.isInitialized() {
		return nil
	}

	return c.client.Close()
}

// Run establishes new connection with the remote host and executes
// the given command.
func (c RemoteClient) Run(command Command) error {
	return c.RunCtx(context.Background(), command)
}

// RunCtx establishes new connection with the remote host and executes
// the given command.
func (c RemoteClient) RunCtx(ctx context.Context, command Command) error {
	// Ensure SSH client is initialized.
	if c.client == nil {
		err := c.initClient(ctx)
		if err != nil {
			c.mux.Unlock()
			return err
		}
	}

	// Initiate new SSH session.
	c.mux.Lock()
	session, err := c.client.NewSession()
	if err != nil {
		c.mux.Unlock()
		return fmt.Errorf("create session for %q: %v", c.Endpoint(), err)
	}
	defer session.Close()
	c.mux.Unlock()

	// Prepare command.
	session.Stdin = c.stdin
	session.Stdout = c.stdout
	session.Stderr = c.stderr

	for k, v := range command.envs {
		err := session.Setenv(k, v)
		if err != nil {
			return fmt.Errorf("set env variable %q: %v", k, err)
		}
	}

	// Run the command.
	cmd := command.command
	if len(command.args) > 0 {
		cmd = fmt.Sprintf("%s %s", command.command, strings.Join(command.args, " "))
	}

	if c.sudo {
		cmd = fmt.Sprintf("sudo --preserve-env %s", cmd)
	}

	return session.Run(cmd)
}

func (c *RemoteClient) initClient(ctx context.Context) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	config := &ssh.ClientConfig{}
	config.User = c.user

	// Inherit command timeout from context.
	timeout, ok := ctx.Deadline()
	if ok {
		config.Timeout = time.Until(timeout)
	}

	// Read public key and set it as valid authorization keys.
	if c.publicKeyPath != "" {
		file, err := os.ReadFile(c.publicKeyPath)
		if err != nil {
			return fmt.Errorf("read public key: %v", err)
		}

		publicKey, _, _, _, err := ssh.ParseAuthorizedKey(file)
		if err != nil {
			return fmt.Errorf("parse public key: %v", err)
		}

		config.HostKeyCallback = ssh.FixedHostKey(publicKey)
	} else {
		config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	// Read private key from the file and set it as a signer for
	// public keys.
	if c.privateKeyPath != "" {
		file, err := os.ReadFile(c.privateKeyPath)
		if err != nil {
			return fmt.Errorf("read private key: %v", err)
		}

		privateKey, err := ssh.ParsePrivateKey(file)
		if err != nil {
			return fmt.Errorf("parse private key: %v", err)
		}

		config.Auth = append(config.Auth, ssh.PublicKeys(privateKey))
	}

	// Connect to the host and store reference to the initialized client.
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", c.host, c.port), config)
	if err != nil {
		return fmt.Errorf("dial %q: %v", c.Endpoint(), err)
	}

	c.client = client
	c.initialized = true

	return nil
}

func (c RemoteClient) isInitialized() bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.initialized
}

// Run is a shorthand for running the given command locally.
func Run(command Command) error {
	return RunCtx(context.Background(), command)
}

// RunCtx is a shorthand for running the given command locally.
func RunCtx(ctx context.Context, command Command) error {
	c := NewLocalClient()
	c.SetStdout(os.Stdout)
	c.SetStderr(os.Stderr)

	return c.RunCtx(ctx, command)
}
