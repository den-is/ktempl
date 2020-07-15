package render

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/den-is/ktempl/pkg/kubernetes"
	"github.com/den-is/ktempl/pkg/logging"
	"github.com/den-is/ktempl/pkg/validation"
	"github.com/spf13/viper"
)

type TemplData struct {
	Nodes  *[]kubernetes.Node
	Values *map[string]interface{}
}

// Receives path to a template, Node data and writes rendered file to output destination path
func ProduceOutput(templPath string, inputData *TemplData, outputDst string) error {

	// prepare template and rendered data
	templName := path.Base(templPath)
	t, templInitErr := template.New(templName).Funcs(sprig.TxtFuncMap()).Option("missingkey=error").ParseFiles(templPath)
	if templInitErr != nil {
		logging.LogWithFields(
			logging.Fields{
				"component": "render",
			}, "error", templInitErr)
	}
	templRenderedBuf := new(bytes.Buffer)
	if templExecErr := t.Execute(templRenderedBuf, *inputData); templExecErr != nil {
		logging.LogWithFields(
			logging.Fields{
				"component": "render",
			}, "error", templExecErr)
		return templExecErr
	}

	outputFilePermsVal := viper.GetUint32("permissions")

	if outputDst == "" {
		// write to stdout if output file path not provided
		fmt.Fprint(os.Stdout, templRenderedBuf.String())
		return nil

	} else if err := validation.CheckFileExists(outputDst); err != nil && outputDst != "" {
		// if file not exists and not stdout -> create file

		outpufFile, err := os.Create(outputDst)
		if err == nil {

			// change permissions of a new file
			if err := outpufFile.Chmod(os.FileMode(outputFilePermsVal)); err != nil {
				logging.LogWithFields(
					logging.Fields{
						"component": "render",
					}, "error", "Failed setting permissions on a new file:", outpufFile, err)
			}

			// make sure to close file after writing it
			defer outpufFile.Close()

			writeErr := ioutil.WriteFile(outputDst, templRenderedBuf.Bytes(), os.FileMode(outputFilePermsVal))
			if writeErr != nil {
				logging.LogWithFields(
					logging.Fields{
						"component": "render",
					}, "error", "Was not able to write new file:", outputDst, err)
				return writeErr

			}

			fmt.Println("Successfully wrote contents in", outputDst)
			logging.LogWithFields(
				logging.Fields{
					"component": "render",
				}, "warn", "Succesfully wrote file ", outputDst)
			return nil

		}

		logging.LogWithFields(
			logging.Fields{
				"component": "render",
			}, "error", "Was not able to create output file", outputDst)
		return err

	} else if outputDst != "" {
		// if file exists and not stdout -> compare contents and [over]write

		// open file
		existingFileBS, err := ioutil.ReadFile(outputDst)
		if err != nil {
			fmt.Println("Unable to read file at destination ", outputDst, err)
			logging.LogWithFields(
				logging.Fields{
					"component": "render",
				}, "error", "Unable to read file at destination ", outputDst, err)
			return err

			// compare data at destination with rendered
		} else if string(existingFileBS) != templRenderedBuf.String() {
			// if datas differ. write rendered data into destination
			writeErr := ioutil.WriteFile(outputDst, templRenderedBuf.Bytes(), os.FileMode(outputFilePermsVal))
			if writeErr != nil {
				// log error if write to existing file fails
				logging.LogWithFields(
					logging.Fields{
						"component": "render",
					}, "error", "Was not able to write to existing output file", outputDst, writeErr)
				return writeErr

			}
			// log success if write to existing file succeeds
			logging.LogWithFields(
				logging.Fields{
					"component": "render",
				}, "warn", "Successfully updated file", outputDst)
			return nil

		} else {
			fmt.Println("Contents did not change")
			logging.LogWithFields(
				logging.Fields{
					"component": "render",
				}, "info", "Contents did not change")
			return errors.New("Contents did not change")
		}
	}

	return nil

}
