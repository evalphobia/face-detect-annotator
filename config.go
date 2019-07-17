package fda

import (
	"fmt"
	"os"
	"path"
	"strconv"
)

const (
	keyConfigEngineAll = "FDA_ENGINE_ALL"
	// Basics
	keyConfigEngineAzureVision  = "FDA_ENGINE_AZURE"
	keyConfigEngineGoogleVision = "FDA_ENGINE_GOOGLE"
	keyConfigEngineRekognition  = "FDA_ENGINE_REKOGNITION"
	keyConfigEnginePigo         = "FDA_ENGINE_PIGO"

	// Extentions
	keyConfigEngineDlib       = "FDA_ENGINE_DLIB"
	keyConfigEngineOpenCV     = "FDA_ENGINE_OPENCV"
	keyConfigEngineTensorFlow = "FDA_ENGINE_TF"
)

const (
	keyConfigAzureRegion          = "FDA_AZURE_REGION"
	keyConfigAzureSubscriptionKey = "FDA_AZURE_SUBSCRIPTION_KEY"
	defaultAzureRegion            = "eastus"

	keyConfigPigoCascadeFilePath = "FDA_PIGO_CASCADE_FILE"
	defaultPigoCascadeFilePath   = "models/facefinder"

	keyConfigDlibModelDir = "FDA_DLIB_MODEL_DIR"
	defaultDlibModelDir   = "models"

	keyConfigOpenCVCascadeFilePath = "FDA_OPENCV_CASCADE_FILE"
	defaultOpenCVCascadeFilePath   = "models/opencv.xml"

	keyConfigTensorFlowModelFilePath = "FDA_TF_MODEL_FILE"
	defaultTensorFlowModelFilePath   = "models/tensorflow.pb"
)

type Config struct {
	InputPath  string
	OutputPath string

	UseEngineAzureVision  bool
	UseEngineGoogleVision bool
	UseEngineRekognition  bool
	UseEnginePigo         bool
	UseEngineDlib         bool
	UseEngineOpenCV       bool
	UseEngineTensorFlow   bool

	AzureRegion           string
	AzureSubscriptionKey  string
	PigoCascadeFilePath   string
	DlibModelDir          string
	OpneCVCascadeFilePath string
	TensorFlowModelPath   string
}

func NewConfig() Config {
	useAzure, _ := strconv.ParseBool(os.Getenv(keyConfigEngineAzureVision))
	useGoogle, _ := strconv.ParseBool(os.Getenv(keyConfigEngineGoogleVision))
	useRekognition, _ := strconv.ParseBool(os.Getenv(keyConfigEngineRekognition))
	usePigo, _ := strconv.ParseBool(os.Getenv(keyConfigEnginePigo))
	useDlib, _ := strconv.ParseBool(os.Getenv(keyConfigEngineDlib))
	useOpenCV, _ := strconv.ParseBool(os.Getenv(keyConfigEngineOpenCV))
	useTF, _ := strconv.ParseBool(os.Getenv(keyConfigEngineTensorFlow))

	useAll, _ := strconv.ParseBool(os.Getenv(keyConfigEngineAll))
	if useAll {
		useAzure = true
		useGoogle = true
		useRekognition = true
		usePigo = true
		useDlib = true
		useOpenCV = true
		useTF = true
	}

	return Config{
		UseEngineAzureVision:  useAzure,
		UseEngineGoogleVision: useGoogle,
		UseEngineRekognition:  useRekognition,
		UseEnginePigo:         usePigo,
		UseEngineDlib:         useDlib,
		UseEngineOpenCV:       useOpenCV,
		UseEngineTensorFlow:   useTF,
		PigoCascadeFilePath:   os.Getenv(keyConfigPigoCascadeFilePath),
		OpneCVCascadeFilePath: os.Getenv(keyConfigOpenCVCascadeFilePath),
		DlibModelDir:          os.Getenv(keyConfigDlibModelDir),
		TensorFlowModelPath:   os.Getenv(keyConfigTensorFlowModelFilePath),
		AzureRegion:           os.Getenv(keyConfigAzureRegion),
		AzureSubscriptionKey:  os.Getenv(keyConfigAzureSubscriptionKey),
	}
}

func (c *Config) setInputPath(s string) {
	c.InputPath = s
}

func (c *Config) setOutputPath(s string) {
	c.OutputPath = s
}

func (c *Config) setUseEngineFromName(name string) error {
	switch name {
	case "azure":
		c.UseEngineAzureVision = true
	case "google":
		c.UseEngineGoogleVision = true
	case "rekognition":
		c.UseEngineRekognition = true
	case "pigo":
		c.UseEnginePigo = true
	case "dlib":
		c.UseEngineDlib = true
	case "opencv":
		c.UseEngineOpenCV = true
	case "tensorflow":
		c.UseEngineTensorFlow = true
	default:
		return fmt.Errorf("unknown engine name: [%s]", name)
	}
	return nil
}

func (c Config) isCSVFilePath() bool {
	switch path.Ext(c.InputPath) {
	case ".csv", ".tsv":
		return true
	}
	return false
}

func (c Config) GetAzureRegion() string {
	if c.AzureRegion != "" {
		return c.AzureRegion
	}
	return defaultAzureRegion
}

func (c Config) GetAzureSubscriptionKey() string {
	return c.AzureSubscriptionKey
}

func (c Config) GetPigoCascadeFile() string {
	if c.PigoCascadeFilePath != "" {
		return c.PigoCascadeFilePath
	}
	return defaultPigoCascadeFilePath
}

func (c Config) GetDlibModelDir() string {
	if c.DlibModelDir != "" {
		return c.DlibModelDir
	}
	return defaultDlibModelDir
}

func (c Config) GetOpenCVCascadeFile() string {
	if c.OpneCVCascadeFilePath != "" {
		return c.OpneCVCascadeFilePath
	}
	return defaultOpenCVCascadeFilePath
}

func (c Config) GetTensorFlowModelFile() string {
	if c.TensorFlowModelPath != "" {
		return c.TensorFlowModelPath
	}
	return defaultTensorFlowModelFilePath
}
