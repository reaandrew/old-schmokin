package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var SchmokinFormat = `content_type: %{content_type}\n filename_effective: %{filename_effective}\n ftp_entry_path: %{ftp_entry_path}\n http_code: %{http_code}\n http_connect: %{http_connect}\n local_ip: %{local_ip}\n local_port: %{local_port}\n num_connects: %{num_connects}\n num_redirects: %{num_redirects}\n redirect_url: %{redirect_url}\n remote_ip: %{remote_ip}\n remote_port: %{remote_port}\n size_download: %{size_download}\n size_header: %{size_header}\n size_request: %{size_request}\n size_upload: %{size_upload}\n speed_download: %{speed_download}\n speed_upload: %{speed_upload}\n ssl_verify_result: %{ssl_verify_result}\n time_appconnect: %{time_appconnect}\n time_connect: %{time_connect}\n time_namelookup: %{time_namelookup}\n time_pretransfer: %{time_pretransfer}\n time_redirect: %{time_redirect}\n time_starttransfer: %{time_starttransfer}\n time_total: %{time_total}\n url_effective: %{url_effective}\n`

func run() {
	processCmd := exec.Command("curl")
	//stdout, err := processCmd.StdoutPipe()
	processCmd.Start()
}

type SchmokinResponse struct {
	response string
	payload  string
}

type SchmokinResult struct {
	success bool
}

type SchmokinHttpClient interface {
	execute(args []string) SchmokinResponse
}

type CurlHttpClient struct {
	args []string
}

func (instance CurlHttpClient) execute(args []string) SchmokinResponse {
	process := "curl"

	executeArgs := append(args, instance.args...)

	var output []byte
	var err error

	if output, err = exec.Command(process, executeArgs...).CombinedOutput(); err != nil {
		exitError := err.(*exec.ExitError)
		fmt.Println(string(exitError.Stderr))
		os.Exit(1)
	}

	payloadData, _ := ioutil.ReadFile("schmokin-response")

	return SchmokinResponse{
		payload:  string(payloadData),
		response: string(output),
	}
}

func CreateCurlHttpClient() CurlHttpClient {
	baseArgs := []string{
		"-v",
		"-s",
		fmt.Sprintf("-w '%s'", SchmokinFormat),
		"-o",
		"schmokin-response",
	}
	return CurlHttpClient{
		args: baseArgs,
	}
}

type SchmokinApp struct {
	httpClient SchmokinHttpClient
	target     string
}

func (instance SchmokinApp) schmoke(args []string) SchmokinResult {

	result := instance.httpClient.execute([]string{args[0]})

	success := true
	current := 0

	for current < len(args) {
		switch args[current] {
		case "--status":
			reg, _ := regexp.Compile(`http_code:\s([\d]+)`)
			result_slice := reg.FindAllStringSubmatch(result.response, -1)
			if len(result_slice) == 1 && len(result_slice[0]) == 2 {
				instance.target = result_slice[0][1]
			}
		case "--filename_effective", "--ftp_entry_path", "--http_code", "--http_connect", "--local_ip", "--local_port", "--num_connects", "--num_redirects", "--redirect_url", "--remote_ip", "--remote_port", "--size_download", "--size_header", "--size_request", "--size_upload", "--speed_download", "--speed_upload", "--ssl_verify_result", "--time_appconnect", "--time_connect", "--time_namelookup", "--time_pretransfer", "--time_redirect", "--time_starttransfer", "--time_total", "--url_effective":
			fmt.Println(fmt.Sprintf("arg = %s", args[current]))
			reg, _ := regexp.Compile(fmt.Sprintf(`%s:\s([\d]+)`, args[current]))
			result_slice := reg.FindAllStringSubmatch(result.response, -1)
			if len(result_slice) == 1 && len(result_slice[0]) == 2 {
				instance.target = result_slice[0][1]
			}
		case "--eq":
			if len(args) < current+2 {
				err := fmt.Errorf("Must supply value to compare against --eq")
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			var expected = args[current+1]
			success = success && (expected == instance.target)
			current += 1
		case "--ne":
			if len(args) < current+2 {
				err := fmt.Errorf("Must supply value to compare against --ne")
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			var expected = args[current+1]
			success = success && (expected != instance.target)
			current += 1
		case "--gt":
			if len(args) < current+2 {
				err := fmt.Errorf("Must supply value to compare against --gt")
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			expected, err := strconv.Atoi(args[current+1])
			if err != nil {
				err = fmt.Errorf("Argument must be a integer for the expected")
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			actual, err := strconv.Atoi(instance.target)
			if err != nil {
				err = fmt.Errorf("Argument must be a integer for the actual")
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			success = success && (actual > expected)
			current += 1
		case "--gte":
			if len(args) < current+2 {
				err := fmt.Errorf("Must supply value to compare against --gte")
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			expected, err := strconv.Atoi(args[current+1])
			if err != nil {
				err = fmt.Errorf("Argument must be a integer for the expected")
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			actual, err := strconv.Atoi(instance.target)
			if err != nil {
				err = fmt.Errorf("Argument must be a integer for the actual")
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			success = success && (actual >= expected)
			current += 1
		case "--lt":
			if len(args) < current+2 {
				err := fmt.Errorf("Must supply value to compare against --lt")
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			expected, err := strconv.Atoi(args[current+1])
			if err != nil {
				err = fmt.Errorf("Argument must be a integer for the expected")
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			actual, err := strconv.Atoi(instance.target)
			if err != nil {
				err = fmt.Errorf("Argument must be a integer for the actual")
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			success = success && (actual < expected)
			current += 1
		case "--lte":
			if len(args) < current+2 {
				err := fmt.Errorf("Must supply value to compare against --lte")
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			expected, err := strconv.Atoi(args[current+1])
			if err != nil {
				err = fmt.Errorf("Argument must be a integer for the expected")
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			actual, err := strconv.Atoi(instance.target)
			if err != nil {
				err = fmt.Errorf("Argument must be a integer for the actual")
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			success = success && (actual <= expected)
			current += 1
		case "--co":
			//TODO: Use --co with other parameters
			if len(args) < current+2 {
				err := fmt.Errorf("Must supply value to compare against --co")
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			var expected = args[current+1]
			success = success && strings.Contains(result.payload, expected)
			current += 1
		case "--resp-header":
			if len(args) < current+2 {
				err := fmt.Errorf("Must supply value to compare against --req-header")
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
			regex := fmt.Sprintf(`(?i)<\s%s:\s([^\n\r]+)`, args[current+1])
			reg, _ := regexp.Compile(regex)
			result_slice := reg.FindAllStringSubmatch(result.response, -1)

			if len(result_slice) == 1 && len(result_slice[0]) == 2 {
				instance.target = result_slice[0][1]
			}
			current += 1
		default:
			if current > 0 {
				panic(fmt.Sprintf("Unknown Arg: %v", args[current]))
			}
		}

		current += 1
	}

	return SchmokinResult{
		success: success,
	}
}

func CreateSchmokinApp(httpClient SchmokinHttpClient) SchmokinApp {
	return SchmokinApp{
		httpClient: httpClient,
	}
}

func main() {

}
