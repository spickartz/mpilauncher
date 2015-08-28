package main

import (
	"fmt"
	"github.com/cheggaaa/pb"
	"github.com/olekukonko/tablewriter"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func run_cmd(params map[string]string) map[string]float64 {
	// Start command and take the time
	start_time := time.Now()
	cmd := exec.Command(params["cmd"], strings.Fields(params["args"])...)
	output, err := cmd.Output()
	duration := float64(time.Since(start_time)) / (1000 * 1000 * 1000)

	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: Could not execute ", cmd, params["args"], err)
		os.Exit(1)
	}
	var result = map[string]float64{
		"overall": duration,
		"inner":   math.NaN(),
		"start":   math.NaN(),
	}

	if params["time_string"] == "" {
		return result
	}

	// Parse output
	lines := strings.Split(string(output), "\n")
	r := regexp.MustCompile(params["time_string"])
	for _, line := range lines {
		if match := r.FindStringSubmatch(line); len(match) > 0 {
			runtime, err := strconv.ParseFloat(match[1], 64)
			if err == nil {
				result["inner"] = runtime
				result["start"] = result["overall"] - runtime
			}
		}
	}

	return result
}

func aggregate_results(exec_results []map[string]float64) []string {
	overall := 0.0
	inner := 0.0
	start := 0.0

	// Calculate sum
	for _, cur_result := range exec_results {
		overall += cur_result["overall"]
		inner += cur_result["inner"]
		start += cur_result["start"]
	}

	// Divide by iterations
	overall /= float64(len(exec_results))
	inner /= float64(len(exec_results))
	start /= float64(len(exec_results))

	var result = []string{
		strconv.FormatFloat(overall, 'f', 2, 64),
		strconv.FormatFloat(inner, 'f', 2, 64),
		strconv.FormatFloat(start, 'f', 2, 64),
	}
	return result
}

func main() {

	var commands = map[string]map[string]string{
		"SP.C.64 (sp/suspend_ib_integration)": {
			"cmd":         "./ib_integration.sh",
			"iterations":  "10",
			"args":        "64 2,3 /home/pickartz/work/benchmarks/NPB/NPB3.3.1/NPB3.3-MPI/bin/sp.B.64",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"SP.C.64 (master)": {
			"cmd":         "./master.sh",
			"iterations":  "10",
			"args":        "64 2,3 /home/pickartz/work/benchmarks/NPB/NPB3.3.1/NPB3.3-MPI/bin/sp.B.64",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
	}

	// Prepare output table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"app", "outer", "inner", "startup"})
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table_content := make([][]string, 0)

	// For all commands
	for app, params := range commands {
		iterations, _ := strconv.ParseInt(params["iterations"], 0, 64)
		cur_results := make([]map[string]float64, 0)

		// Perform requested iterations
		fmt.Printf("STATUS: Executing '%s' ...\n", app)
		progress_bar := pb.StartNew(int(iterations))
		for i := int64(0); i < iterations; i++ {
			cur_results = append(cur_results, run_cmd(params))
			progress_bar.Increment()
		}
		progress_bar.FinishPrint("Done")

		// Aggregate results
		table_row := make([]string, 0)
		table_row = append(table_row, app)
		table_row = append(table_row, aggregate_results(cur_results)...)
		table_content = append(table_content, table_row)
	}

	table.AppendBulk(table_content)
	table.Render()
}
