package cmd

// exportKubeconfigCmd represents the exportKubeconfig command
// var exportKubeconfigCmd = &cobra.Command{
// 	Use:   "kubeconfig",
// 	Short: "Export cluster kubeconfig file",
// 	Long:  `Command export kubeconfig prints content of the kubeconfig file.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		err := exportKubeconfig()
// 		if err != nil {
// 			fmt.Fprintln(os.Stderr, err)
// 			os.Exit(1)
// 		}
// 	},
// }

// func init() {
// 	exportCmd.AddCommand(exportKubeconfigCmd)

// 	exportKubeconfigCmd.PersistentFlags().StringVar(&env.ClusterName, "cluster", env.DefaultClusterName, "specify the cluster to be used")
// 	exportKubeconfigCmd.PersistentFlags().BoolVarP(&env.Local, "local", "l", false, "use a current directory as the cluster path")

// 	// Auto complete cluster names of active clusters that also contain kubeconfig
// 	// for the flag 'cluster'.
// 	exportKubeconfigCmd.RegisterFlagCompletionFunc("cluster", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

// 		clusterNames, err := GetClusters([]ClusterFilter{IsActive, ContainsKubeconfig})
// 		if err != nil {
// 			return nil, cobra.ShellCompDirectiveNoFileComp
// 		}

// 		return clusterNames, cobra.ShellCompDirectiveNoFileComp
// 	})
// }

// // exportKubeconfig exports (prints) content of the cluster Kubeconfig file.
// func exportKubeconfig() error {

// 	kubeconfigPath := filepath.Join(env.ClusterPath, env.ConstKubeconfigPath)

// 	err := utils.VerifyClusterDir(env.ClusterPath)
// 	if err != nil {
// 		return fmt.Errorf("Cluster '%s' does not exist: %w", env.ClusterName, err)
// 	}

// 	_, err = os.Stat(kubeconfigPath)
// 	if err != nil {
// 		return fmt.Errorf("Kubeconfig for cluster '%s' does not exist: %w", env.ClusterName, err)
// 	}

// 	kubeconfig, err := ioutil.ReadFile(kubeconfigPath)
// 	if err != nil {
// 		return fmt.Errorf("Failed reading Kubeconfig file: %w", err)
// 	}

// 	fmt.Print(string(kubeconfig))

// 	return nil
// }
