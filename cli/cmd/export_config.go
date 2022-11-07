package cmd

// var exportConfigCmd = &cobra.Command{
// 	Use:   "config",
// 	Short: "Export cluster configuration file",
// 	Long: `
// 	Command export config prints content of the cluster configuration file.`,

// 	Run: func(cmd *cobra.Command, args []string) {
// 		err := exportConfig()
// 		if err != nil {
// 			fmt.Fprintln(os.Stderr, err)
// 			os.Exit(1)
// 		}
// 	},
// }

// func init() {
// 	exportCmd.AddCommand(exportConfigCmd)

// 	exportConfigCmd.PersistentFlags().StringVar(&env.ClusterName, "cluster", env.DefaultClusterName, "specify the cluster to be used")
// 	exportConfigCmd.PersistentFlags().BoolVarP(&env.Local, "local", "l", false, "use a current directory as the cluster path")

// 	// Auto complete cluster names of clusters that contain Kubitect config
// 	// for the flag 'cluster'
// 	exportConfigCmd.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

// 		clusterNames, err := GetClusters([]ClusterFilter{ContainsConfig})
// 		if err != nil {
// 			return nil, cobra.ShellCompDirectiveNoFileComp
// 		}

// 		return clusterNames, cobra.ShellCompDirectiveNoFileComp
// 	})
// }

// // exportConfig exports (prints) content of the cluster configuration file.
// func exportConfig() error {

// 	configPath := filepath.Join(env.ClusterPath, env.DefaultClusterConfigPath)

// 	err := utils.VerifyClusterDir(env.ClusterPath)
// 	if err != nil {
// 		return fmt.Errorf("Cluster '%s' does not exist: %w", env.ClusterName, err)
// 	}

// 	_, err = os.Stat(configPath)
// 	if err != nil {
// 		return fmt.Errorf("Cluster configuration for cluster '%s' does not exist: %w", env.ClusterName, err)
// 	}

// 	config, err := ioutil.ReadFile(configPath)
// 	if err != nil {
// 		return fmt.Errorf("Failed reading Kubeconfig file: %w", err)
// 	}

// 	fmt.Print(string(config))

// 	return nil
// }
