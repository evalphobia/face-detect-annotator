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

	// Extentions
	keyConfigEngineDlib       = "FDA_ENGINE_DLIB"
	keyConfigEngineOpenCV     = "FDA_ENGINE_OPENCV"
	keyConfigEngineTensorFlow = "FDA_ENGINE_TF"
)

const (
	keyConfigAzureRegion          = "FDA_AZURE_REGION"
	keyConfigAzureSubscriptionKey = "FDA_AZURE_SUBSCRIPTION_KEY"
	defaultAzureRegion            = "eastus"

	keyConfigDlibModelDir = "FDA_DLIB_MODEL_DIR"
	defaultDlibModelDir   = "models"

	keyConfigOpenCVCascadeFilePath = "FDA_OPENCV_CASCADE_FILE"
	defaultOpenCVCascadeFilePath   = "models/opencv.xml"

	keyConfigTensorFlowModelFilePath = "FDA_TF_MODEL_FILE"
	defaultTensorFlowModelFilePath   = "models/tensorflow.pb"
)

type Config struct {
	InputPath             string
	OutputPath            string
	UseEngineOpenCV       bool
	UseEngineDlib         bool
	UseEngineTensorFlow   bool
	UseEngineRekognition  bool
	UseEngineGoogleVision bool
	UseEngineAzureVision  bool
	OpneCVCascadeFilePath string
	DlibModelDir          string
	TensorFlowModelPath   string
	AzureRegion           string
	AzureSubscriptionKey  string
}

func NewConfig() Config {
	useOpenCV, _ := strconv.ParseBool(os.Getenv(keyConfigEngineOpenCV))
	useDlib, _ := strconv.ParseBool(os.Getenv(keyConfigEngineDlib))
	useTF, _ := strconv.ParseBool(os.Getenv(keyConfigEngineTensorFlow))
	useRekognition, _ := strconv.ParseBool(os.Getenv(keyConfigEngineRekognition))
	useGoogle, _ := strconv.ParseBool(os.Getenv(keyConfigEngineGoogleVision))
	useAzure, _ := strconv.ParseBool(os.Getenv(keyConfigEngineAzureVision))
	useAll, _ := strconv.ParseBool(os.Getenv(keyConfigEngineAll))
	if useAll {
		useOpenCV = true
		useDlib = true
		useTF = true
		useRekognition = true
		useGoogle = true
		useAzure = true
	}

	return Config{
		UseEngineOpenCV:       useOpenCV,
		UseEngineDlib:         useDlib,
		UseEngineTensorFlow:   useTF,
		UseEngineRekognition:  useRekognition,
		UseEngineGoogleVision: useGoogle,
		UseEngineAzureVision:  useAzure,
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
	case "opencv":
		c.UseEngineOpenCV = true
	case "dlib":
		c.UseEngineDlib = true
	case "rekognition":
		c.UseEngineRekognition = true
	case "google":
		c.UseEngineGoogleVision = true
	case "azure":
		c.UseEngineAzureVision = true
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
