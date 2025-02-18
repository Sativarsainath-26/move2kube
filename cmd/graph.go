/*
 *  Copyright IBM Corporation 2022
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/konveyor/move2kube/common"
	graphutils "github.com/konveyor/move2kube/graph"
	"github.com/konveyor/move2kube/types"
	graphtypes "github.com/konveyor/move2kube/types/graph"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type graphFlags struct {
	graphFilePath string
	port          int32
	outputPath    string
}

func graphHandler(flags graphFlags) {
	graphFilePath := filepath.Clean(flags.graphFilePath)
	outputPath := filepath.Clean(flags.outputPath)
	graphFile, err := os.Open(graphFilePath)
	if err != nil {
		logrus.Fatalf("failed to the open the graph file at path %s . Error: %q", graphFilePath, err)
	}
	graph := graphtypes.Graph{}
	if err := json.NewDecoder(graphFile).Decode(&graph); err != nil {
		logrus.Fatalf("failed to decode the json file at path %s . Error: %q", graphFilePath, err)
	}
	nodes, edges := graphutils.GetNodesAndEdges(graph)
	graphutils.DfsUpdatePositions(nodes, edges)
	webGraph := graphtypes.GraphT{Nodes: nodes, Edges: edges}
	if flags.outputPath != "" {
		webBytes, err := json.Marshal(webGraph)
		if err != nil {
			logrus.Fatalf("failed to marshal the processed graph to json. Error: %q", err)
		}
		if err := os.WriteFile(outputPath, webBytes, common.DefaultFilePermission); err != nil {
			logrus.Fatalf("failed to write the processed graph json to a file at path %s . Error: %q", outputPath, err)
		}
		return
	}
	logrus.Fatalf("graph server stopped. Error: %q", graphutils.StartServer(webGraph, flags.port))
}

// GetGraphCommand returns a command to show the graph of all the transformers that were run
func GetGraphCommand() *cobra.Command {
	viper.AutomaticEnv()
	flags := graphFlags{}
	graphCmd := &cobra.Command{
		Use:   "graph [-g path/to/m2k-graph.json]",
		Short: "View the graph generated by transform command.",
		Long: `View the graph generated by transform command. This command starts a server to serve a web UI and display the graph.
	To see the graph, go to http://localhost:8080/ in a browser.
	By default, it will look for the m2k-graph.json file in the current working directory.`,
		Run: func(_ *cobra.Command, __ []string) { graphHandler(flags) },
	}
	graphCmd.Flags().StringVarP(&flags.graphFilePath, "graph", "f", "m2k-graph.json", "Path to a m2k-graph.json file generated by the transform command.")
	graphCmd.Flags().Int32VarP(&flags.port, "port", "p", 8080, "Port to start the server on.")
	graphCmd.Flags().StringVarP(&flags.outputPath, "output", "o", "", "Path where the processed graph json file should be generated. If this flag is used then instead of starting a web server, we will output a file. By default "+types.AppName+" does not output this file.")
	return graphCmd
}
