package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmccormick2001/rqlite-operator/pkg/apis/rqcluster/v1alpha1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var rootCmd = &cobra.Command{
	Use:  "rqo",
	Long: `rqlite-operator CLI`,
}

var (
	home     = os.Getenv("HOME")
	username string
	password string

	// ConfigFile is $HOME/.rqo/config.json per default
	// contains user, password and url of zapper
	configFile string
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	addCommands()

	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "username, faculzapperive if you have a "+home+"/.rqo/config.json file")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password, faculzapperive if you have a "+home+"/.rqo/config.json file")
	rootCmd.PersistentFlags().StringVarP(&configFile, "configFile", "c", home+"/.rqo/config.json", "configuration file, default is "+home+"/.rqo/config.json")

	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))

	log.SetLevel(log.DebugLevel)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

//AddCommands adds child commands to the root command rootCmd.
func addCommands() {
	rootCmd.AddCommand(cmdCreate)
}

var cmdCreate = &cobra.Command{
	Use:   "create <rqlite cluster name>",
	Short: "Create a rqlite cluster",
	Run: func(cmd *cobra.Command, args []string) {
		create(args)
	},
}

// create creates a message in specified topic
func create(name []string) {
	readConfig()
	if len(name) > 0 {
		for i := 0; i < len(name); i++ {
			err := createCR(name[i])
			if err != nil {
				log.Errorf("could not create rqlite cluster %s", name[i])
			} else {
				log.Debugf("rqlite cluster %s Created", name[i])
			}
		}
	}
}

/**
func getClient() *zapper.Client {
	tc, err := zapper.NewClient(zapper.Options{
		URL:      viper.GetString("url"),
		Username: viper.GetString("username"),
		Password: viper.GetString("password"),
		Referer:  "rqo.v0",
	})

	if err != nil {
		log.Fatalf("Error while create new zapper Client: %s", err)
	}

	zapper.DebugLogFunc = log.Debugf
	return tc
}
*/

// readConfig reads config in .rqo/config per default
func readConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
		viper.ReadInConfig() // Find and read the config file
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func createCR(name string) error {
	var err error
	crSpec := v1alpha1.RqclusterSpec{
		Size:          3,
		CpuLimit:      "",
		CpuRequest:    "",
		MemoryLimit:   "",
		MemoryRequest: "",
		StorageClass:  "fast",
		StorageLimit:  "10Mi",
	}
	if crSpec.Size > 4 {
	}

	return err
}
