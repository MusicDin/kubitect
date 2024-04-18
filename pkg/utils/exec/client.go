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
var _ Client = localClient{}
var _ Client = remoteClient{}

type Client interface {
	Run(command string, args ...string) error
	RunCtx(ctx context.Context, command string, args ...string) error
	Close() error
}

// commonClient provides functions that are common for all clients.
type commonClient struct {
	// Streams.
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	// Environment variables.
	envs map[string]string
}

func newCommonClient() commonClient {
	return commonClient{
		envs: map[string]string{
			// Set default PATH variable.
			"PATH": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
		},
	}
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

func (c *commonClient) SetEnv(key string, value string) {
	c.envs[key] = value
}

func (c *commonClient) Envs() map[string]string {
	// Return a copy, to prevent uncontrolled modifications
	// or original map.
	envs := make(map[string]string, len(c.envs))
	for k, v := range c.envs {
		envs[k] = v
	}

	return envs
}

func (c commonClient) Close() error {
	return nil
}

type localClient struct {
	commonClient

	workingDir string
}

// NewLocalClient initializes a client for running local commands.
func NewLocalClient() localClient {
	return localClient{
		commonClient: newCommonClient(),
	}
}

// Set working directory in which commands are executed.
func (c localClient) WithWorkingDir(wd string) localClient {
	c.workingDir = wd
	return c
}

// Run runs command locally.
func (c localClient) Run(command string, args ...string) error {
	return c.RunCtx(context.Background(), command, args...)
}

// RunCtx runs command locally.
func (c localClient) RunCtx(ctx context.Context, command string, args ...string) error {
	command, args = splitOneLineCommand(command, args)

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdin = c.stdin
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr

	env := make([]string, 0, len(c.envs))
	for k, v := range c.envs {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	cmd.Env = env
	cmd.Dir = c.workingDir

	return cmd.Run()
}

// Output runs command locally and returns standard output as slice of bytes.
func (c localClient) Output(command string, args ...string) (stdout []byte, err error) {
	return c.OutputCtx(context.Background(), command, args...)
}

// OutputCtx runs command locally and returns standard output as slice of bytes.
func (c localClient) OutputCtx(ctx context.Context, command string, args ...string) (stdout []byte, err error) {
	return c.buildCommand(ctx, command, args...).Output()
}

func (c localClient) buildCommand(ctx context.Context, command string, args ...string) *exec.Cmd {
	command, args = splitOneLineCommand(command, args)

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdin = c.stdin
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stderr

	env := make([]string, 0, len(c.envs))
	for k, v := range c.envs {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	cmd.Env = env
	cmd.Dir = c.workingDir

	return cmd
}

type remoteClient struct {
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
func NewSSHClient(user string, host string) remoteClient {
	return remoteClient{
		commonClient: newCommonClient(),
		user:         user,
		host:         host,
		port:         "22",
		mux:          &sync.Mutex{},
	}
}

// WithPort sets the host port to given value. By default, port 22 is used.
// Immutable once client is initialized.
func (c remoteClient) WithPort(port uint16) remoteClient {
	if !c.isInitialized() {
		c.port = fmt.Sprint(port)
	}

	return c
}

// WithPrivateKeyFile sets the path to the private key file that is used
// for authentication. Immutable once client is initialized.
func (c remoteClient) WithPrivateKeyFile(privateKeyPath string) remoteClient {
	if !c.isInitialized() {
		c.privateKeyPath = privateKeyPath
	}

	return c
}

// WithPublicKeyFile sets the path to the public key file that is used
// for host verification. If not set, known hosts are ignored.
// Immutable once client is initialized.
func (c remoteClient) WithPublicKeyFile(publicKeyPath string) remoteClient {
	if !c.isInitialized() {
		c.publicKeyPath = publicKeyPath
	}

	return c
}

// WithSuperUser runs the command as super user, effectively prepending
// "sudo" in from of the command. Immutable once client is initialized.
func (c remoteClient) WithSuperUser(sudo bool) remoteClient {
	if !c.isInitialized() {
		c.sudo = sudo
	}

	return c
}

// Endpoint returns SSH endpoint in format "user@host:port".
func (c remoteClient) Endpoint() string {
	return fmt.Sprintf("%s@%s:%s", c.user, c.host, c.port)
}

// Close closes potentially initialized the SSH client.
func (c remoteClient) Close() error {
	if !c.isInitialized() {
		return nil
	}

	return c.client.Close()
}

// Run establishes new connection with the remote host and executes
// the given command.
func (c remoteClient) Run(command string, args ...string) error {
	return c.RunCtx(context.Background(), command)
}

// RunCtx establishes new connection with the remote host and executes
// the given command.
func (c remoteClient) RunCtx(ctx context.Context, command string, args ...string) error {
	command, args = splitOneLineCommand(command, args)

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

	for k, v := range c.envs {
		err := session.Setenv(k, v)
		if err != nil {
			return fmt.Errorf("set env variable %q: %v", k, err)
		}
	}

	// Run the command.
	cmd := command
	if len(args) > 0 {
		cmd = fmt.Sprintf("%s %s", command, strings.Join(args, " "))
	}

	if c.sudo {
		cmd = fmt.Sprintf("sudo --preserve-env %s", cmd)
	}

	return session.Run(cmd)
}

func (c *remoteClient) initClient(ctx context.Context) error {
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

func (c remoteClient) isInitialized() bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.initialized
}

// Run is a shorthand for running the given command locally.
func Run(command string, args ...string) error {
	return RunCtx(context.Background(), command, args...)
}

// RunCtx is a shorthand for running the given command locally.
func RunCtx(ctx context.Context, command string, args ...string) error {
	c := NewLocalClient()
	c.SetStdout(os.Stdout)
	c.SetStderr(os.Stderr)

	return c.RunCtx(ctx, command)
}

// splitOneLineCommand splits the command by spaces when no list of
// arguments is empty. This prevents spaces in commands but allows
// passing commands as a single string as long as input arguments
// do not contain spaces.
func splitOneLineCommand(command string, args []string) (string, []string) {
	if len(args) == 0 {
		split := strings.Split(command, " ")
		command = split[0]
		if len(split) > 1 {
			args = split[1:]
		}
	}

	return command, args
}
