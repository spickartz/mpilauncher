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
	var overall = map[string]float64{
		"avg": 0.0,
		"var": 0.0,
	}
	var inner = map[string]float64{
		"avg": 0.0,
		"var": 0.0,
	}
	var start = map[string]float64{
		"avg": 0.0,
		"var": 0.0,
	}

	// Calculate sum
	for _, cur_result := range exec_results {
		overall["avg"] += cur_result["overall"]
		inner["avg"] += cur_result["inner"]
		start["avg"] += cur_result["start"]
	}

	// Divide by iterations
	overall["avg"] /= float64(len(exec_results))
	inner["avg"] /= float64(len(exec_results))
	start["avg"] /= float64(len(exec_results))

	for _, cur_result := range exec_results {
		overall["var"] += math.Pow(overall["avg"]-cur_result["overall"], 2)
		inner["var"] += math.Pow(inner["avg"]-cur_result["inner"], 2)
		start["var"] += math.Pow(start["avg"]-cur_result["start"], 2)
	}

	overall["stdev"] = math.Sqrt(overall["var"])
	inner["stdev"] = math.Sqrt(inner["var"])
	start["stdev"] = math.Sqrt(start["var"])

	var result = []string{
		strconv.FormatFloat(overall["avg"], 'f', 2, 64),
		strconv.FormatFloat(overall["var"], 'f', 6, 64),
		strconv.FormatFloat(overall["stdev"], 'f', 6, 64),
		strconv.FormatFloat(inner["avg"], 'f', 2, 64),
		strconv.FormatFloat(inner["var"], 'f', 6, 64),
		strconv.FormatFloat(inner["stdev"], 'f', 6, 64),
		strconv.FormatFloat(start["avg"], 'f', 2, 64),
		strconv.FormatFloat(start["var"], 'f', 6, 64),
		strconv.FormatFloat(start["stdev"], 'f', 6, 64),
	}
	return result
}

const runtime = 600.00

func calc_iterations(exec_result map[string]float64) int64 {
	app_run, _ := exec_result["overall"]

	iterations := (runtime / app_run)
	if iterations < 1 {
		iterations = 1
	}
	return int64(iterations)
}

func main() {

	var commands = map[string]map[string]string{
		"LU-MZ.C.4": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "4 lu-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"LU-MZ.C.9": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "9 lu-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"LU-MZ.C.16": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "16 lu-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"LU-MZ.C.25": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "25 lu-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"LU-MZ.C.36": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "36 lu-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"LU-MZ.C.39": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "39 lu-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"SP-MZ.C.4": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "4 sp-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"SP-MZ.C.9": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "9 sp-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"SP-MZ.C.16": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "16 sp-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"SP-MZ.C.25": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "25 sp-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"SP-MZ.C.36": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "36 sp-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"SP-MZ.C.39": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "39 sp-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"BT-MZ.C.4": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "4 bt-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"BT-MZ.C.9": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "9 bt-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"BT-MZ.C.16": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "16 bt-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"BT-MZ.C.25": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "25 bt-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"BT-MZ.C.36": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "36 bt-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
		"BT-MZ.C.39": {
			"cmd":         "./npb_mz_launcher.sh",
			"iterations":  "10",
			"runtime":     "26",
			"args":        "39 bt-mz.C.1",
			"time_string": `Time in seconds\s=\s*(\d+\.\d+)`,
		},
	}

	// Prepare output table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.SetHeader([]string{"app", "outer", "var", "stdev", "inner", "var", "stdev", "startup", "var", "stdev"})
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table_content := make([][]string, 0)

	// For all commands
	for app, params := range commands {
		// create result array
		cur_results := make([]map[string]float64, 0)

		//	iterations, _ := strconv.ParseInt(params["iterations"], 0, 64)
		// Calc iterations based on one execution
		cur_results = append(cur_results, run_cmd(params))
		iterations := calc_iterations(cur_results[0])

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
	return
	table.AppendBulk(table_content)
	table.Render()
}
