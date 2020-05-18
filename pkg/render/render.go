package render

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"text/template"

	"github.com/den-is/ktempl/pkg/kubernetes"
	"github.com/den-is/ktempl/pkg/logging"
	"github.com/den-is/ktempl/pkg/validation"
	"github.com/spf13/viper"
)

type TemplData struct {
	Nodes  *[]kubernetes.Node
	Values *map[string]string
}

// Receives path to a template, Node data and writes rendered file to output destination path
func RenderOutput(tpl_path string, nodedata *TemplData, output_dst string) error {

	// prepare template and rendered data
	tpl_name := path.Base(tpl_path)
	t, templ_init_err := template.New(tpl_name).Option("missingkey=error").ParseFiles(tpl_path)
	if templ_init_err != nil {
		logging.LogWithFields(
			logging.Fields{
				"component": "render",
			}, "error", templ_init_err)
	}
	rendered_tpl_bf := new(bytes.Buffer)
	if templ_exec_err := t.Execute(rendered_tpl_bf, *nodedata); templ_exec_err != nil {
		logging.LogWithFields(
			logging.Fields{
				"component": "render",
			}, "error", templ_exec_err)
		return templ_exec_err
	}

	f_perms_val := viper.GetUint32("permissions")

	if output_dst == "" {
		// write to stdout if output file path not provided
		fmt.Fprint(os.Stdout, rendered_tpl_bf.String())
		return nil

	} else if err := validation.CheckFileExists(output_dst); err != nil && output_dst != "" {
		// if file not exists and not stdout -> create file

		out_f, err := os.Create(output_dst)
		if err == nil {

			// change permissions of a new file
			if err := out_f.Chmod(os.FileMode(f_perms_val)); err != nil {
				logging.LogWithFields(
					logging.Fields{
						"component": "render",
					}, "error", "Failed setting permissions on a new file:", out_f, err)
			}

			// make sure to close file after writing it
			defer out_f.Close()

			write_err := ioutil.WriteFile(output_dst, rendered_tpl_bf.Bytes(), os.FileMode(f_perms_val))
			if write_err != nil {
				logging.LogWithFields(
					logging.Fields{
						"component": "render",
					}, "error", "Was not able to write new file:", output_dst, err)
				return write_err

			} else {
				fmt.Println("Successfully wrote contents in", output_dst)
				logging.LogWithFields(
					logging.Fields{
						"component": "render",
					}, "info", "Succesfully wrote file ", output_dst)
				return nil

			}
		} else {
			logging.LogWithFields(
				logging.Fields{
					"component": "render",
				}, "error", "Was not able to create output file", output_dst)
			return err
		}

	} else if output_dst != "" {
		// if file exists and not stdout -> compare contents and [over]write

		// open file
		dst_bs, err := ioutil.ReadFile(output_dst)
		if err != nil {
			fmt.Println("Unable to read file at destination ", output_dst, err)
			logging.LogWithFields(
				logging.Fields{
					"component": "render",
				}, "error", "Unable to read file at destination ", output_dst, err)
			return err

			// compare data at destination with rendered
		} else if string(dst_bs) != rendered_tpl_bf.String() {
			// if datas differ. write rendered data into destination
			write_err := ioutil.WriteFile(output_dst, rendered_tpl_bf.Bytes(), os.FileMode(f_perms_val))
			if write_err != nil {
				// log error if write to existing file fails
				logging.LogWithFields(
					logging.Fields{
						"component": "render",
					}, "error", "Was not able to write to existing output file", output_dst, write_err)
				return write_err

			} else {
				// log success if write to existing file succeeds
				logging.LogWithFields(
					logging.Fields{
						"component": "render",
					}, "info", "Successfully updated file", output_dst)
				return nil
			}
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
