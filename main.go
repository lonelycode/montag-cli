package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
	"github.com/lonelycode/montag-cli/client"
	"github.com/lonelycode/montag-cli/interfaces"
	"github.com/lonelycode/montag-cli/models"
	aifuncRunner "github.com/lonelycode/montag-cli/scriptExtensions/aifunc_runner"
	"github.com/lonelycode/montag-cli/scriptExtensions/contextInjector"
	"github.com/lonelycode/montag-cli/scriptExtensions/dummyFunc"
	"github.com/lonelycode/montag-cli/scriptExtensions/kvstore"
	"github.com/lonelycode/montag-cli/scriptExtensions/readable"
	scripthttpcaller "github.com/lonelycode/montag-cli/scriptExtensions/script_httpcaller"
	scriptBytesHttpCaller "github.com/lonelycode/montag-cli/scriptExtensions/script_httpcaller_bytes"
	secretGetter "github.com/lonelycode/montag-cli/scriptExtensions/secretGetter"
	snippetStore "github.com/lonelycode/montag-cli/scriptExtensions/snippets"
	vectorlookup "github.com/lonelycode/montag-cli/scriptExtensions/vector_lookup"
	"github.com/n3integration/classifier/naive"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "run montag scripts locally, with resources provided by a montag server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "data",
						Value: "",
						Usage: "file to use as input data",
					},
					&cli.StringFlag{
						Name:  "message",
						Value: "",
						Usage: "user message (simulated interaction with bot)",
					},
					&cli.BoolFlag{
						Name:  "passive",
						Value: false,
						Usage: "set to true to run in passive mode (no user mention)",
					},
				},
				Action: func(cCtx *cli.Context) error {
					apiClient := client.NewClient(cCtx.String("key"), cCtx.String("server"))

					inputs := map[string]interface{}{
						"message":    "",
						"history":    []string{},
						"context":    []string{},
						"userID":     1,
						"is_passive": cCtx.Bool("passive"),
					}

					if cCtx.String("data") != "" {
						data := openFile(cCtx.String("data"))
						err := json.Unmarshal(data, &inputs)
						if err != nil {
							return fmt.Errorf("failed to parse input data: %s", err)
						}
					}

					if cCtx.String("message") != "" {
						fmt.Println("setting message in inputs, overrides any message key in input data")
						inputs["message"] = cCtx.String("message")
					}

					db := getDB()

					so, err := runsScript(cCtx.Args().Get(0), inputs, apiClient, db)
					if err != nil {
						return err
					}

					fmt.Printf("SCRIPT OUTPUT:\n - Output Vars: %v\n - Forward Output: %v\n - Return Override: %v\n\n", so.Outputs, so.Response, so.ReturnQuery)
					return nil
				},
			},
			{
				Name:  "classify",
				Usage: "run a query against a dataset for a classifier pre-upload",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Value:    "",
						Usage:    "the csv file to ingest (no header row)",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "prompt",
						Value:    "",
						Usage:    "the prompt to classify",
						Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					fileBytes := openFile(cCtx.String("file"))
					csvReader := csv.NewReader(bytes.NewReader(fileBytes))
					classifier := naive.New()
					i := 0
					for {
						rec, err := csvReader.Read()
						if err == io.EOF {
							break
						}
						if err != nil {
							log.Fatal(err)
							break
						}

						if len(rec) != 2 {
							log.Fatal("CSV file must have two columns")
							break
						}

						classifier.TrainString(rec[0], rec[1])
						i++
					}

					fmt.Println("Training complete, trained on", i, "records")

					probs, best := classifier.Probabilities(cCtx.String("prompt"))

					fmt.Printf("Classification: %s\n", best)
					fmt.Printf("Probabilities: %v\n", probs)
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "server",
				Value:    "https://montag.example.com",
				Usage:    "server to use for montag resources",
				EnvVars:  []string{"MONTAG_SERVER"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "key",
				Value:    "YOURAPIKEY",
				Usage:    "API key to validate against the server",
				EnvVars:  []string{"MONTAG_KEY"},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "storage",
				Value:   "./montag.sqlite",
				Usage:   "database to use for value storage",
				EnvVars: []string{"MONTAG_DB"},
			},
		},
		Name:  "montag-cli",
		Usage: "run montag scripts locally, with resources provided by a montag server",
		Action: func(cCtx *cli.Context) error {
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getDB() *gorm.DB {
	var d *gorm.DB
	binaryPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	d = GetSqlite(filepath.Dir(binaryPath))

	d.AutoMigrate(&models.ScriptKVStore{})

	return d
}

func GetSqlite(basePath string) *gorm.DB {
	// Open a database connection
	dbPath := basePath + "/montag-cli.db"
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	return db
}

func openFile(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func runsScript(scriptName string, inputs map[string]interface{}, apiClient *client.Client, db *gorm.DB) (*models.ScriptOutput, error) {
	if scriptName == "" {
		return nil, fmt.Errorf("a tengo script filename is required")
	}

	prompt := interfaces.NewPrompt("montag-cli", "montag-cli", nil, nil, nil, "montag-cli", "montag-cli")

	script := openFile(scriptName)

	s := tengo.NewScript(script)
	s.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))
	err := s.Add("montagRun", aifuncRunner.NewAIFuncRunner(apiClient))
	if err != nil {
		return nil, err
	}

	err = s.Add("montagRunAsync", dummyFunc.NewDummyFunc("montagRunAsync", models.DummyMultiCallerResponse()))
	if err != nil {
		return nil, err
	}

	err = s.Add("montagSendMessage", dummyFunc.NewDummyFunc("montagSendMessage", nil))
	if err != nil {
		return nil, err
	}

	err = s.Add("montagMakeHttpRequest", scripthttpcaller.NewHttpCaller())
	if err != nil {
		return nil, err
	}

	err = s.Add("montagHistoryLength", 0)
	if err != nil {
		return nil, fmt.Errorf("failed to add user history: %s", err)
	}

	mRunnerStub := aifuncRunner.NewAIFuncRunner(apiClient)
	mRunnerStub.Managed = true
	err = s.Add("montagManagedRun", mRunnerStub)
	if err != nil {
		return nil, err
	}

	err = s.Add("montagMakeHttpRequestForBytes", scriptBytesHttpCaller.NewBytesHttpCaller())
	if err != nil {
		return nil, err
	}

	err = s.Add("montagGetReadableURL", readable.NewReadable())
	if err != nil {
		return nil, err
	}

	err = s.Add("montagVectorSearch", vectorlookup.NewVectorLookup(apiClient))
	if err != nil {
		return nil, err
	}

	err = s.Add("montagAddToHistory", dummyFunc.NewDummyFunc("montagAddToHistory", &tengo.Int{Value: 1}))
	if err != nil {
		return nil, err
	}

	// err = s.Add("montagAddToContext", dummyFunc.NewDummyFunc("montagAddToContext", &tengo.Int{Value: 1}))
	// if err != nil {
	// 	return nil, err
	// }

	err = s.Add("montagAddToContext", contextInjector.NewContextInjector(prompt))
	if err != nil {
		return nil, err
	}

	err = s.Add("montagGetSecret", secretGetter.NewSecretGetter())
	if err != nil {
		return nil, err
	}

	err = s.Add("montagKV", kvstore.NewKVStore(db, 1))
	if err != nil {
		return nil, err
	}

	err = s.Add("montagGetSnippet", snippetStore.NewSnippetStore(apiClient))
	if err != nil {
		return nil, err
	}

	msg, ok := inputs["message"]
	if !ok {
		msg = "undefined"
	}

	err = s.Add("montagUserMessage", msg)
	if err != nil {
		return nil, fmt.Errorf("failed to add user message: %s", err)
	}

	err = s.Add("montagIsMsgPassive", inputs["is_passive"])
	if err != nil {
		return nil, fmt.Errorf("failed to add montagIsMsgPassive: %s", err)
	}

	tHistory := &tengo.Array{}
	tString := &tengo.String{Value: "User: DUMMY INPUT"}
	tHistory.Value = append(tHistory.Value, tString)

	err = s.Add("montagUserHistory", tHistory)
	if err != nil {
		return nil, fmt.Errorf("failed to add user history: %s", err)
	}

	tResources := &tengo.Array{}
	err = s.Add("montagResources", tResources)
	if err != nil {
		return nil, fmt.Errorf("failed to add resources: %s", err)
	}

	// add the remainder
	for k, v := range inputs {
		if k == "message" || k == "history" || k == "context" {
			continue
		}
		key := fmt.Sprintf("montag_%s", k)
		err = s.Add(key, v)
		if err != nil {
			return nil, fmt.Errorf("input remainder import error: %s (%s)", err, key)
		}

	}

	compiled, err := s.Run()
	if err != nil {
		return nil, fmt.Errorf("compile-time error: %s", err)
	}

	so := &models.ScriptOutput{}
	response := compiled.Get("montagResponse")
	if response != nil {
		//fmt.Println("script has response of: ", response.String())
		so.Response = response.String()
	}

	outputs := compiled.Get("montagOutputs")
	if outputs != nil {
		// fmt.Println("script has outputs of: ", outputs.Map())
		so.Outputs = outputs.Map()
	}

	returnQuery := compiled.Get("montagOverride")
	if returnQuery != nil {
		// fmt.Println("script has reply of: ", returnQuery.String())
		so.ReturnQuery = returnQuery.String()
	}

	return so, nil
}
