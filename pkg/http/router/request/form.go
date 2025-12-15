package request

import (
	"bufio"
	"fmt"
	"net/http"
)

const fileFormLimit = 1024 << 20

func (r *request) parseMultipartForm(req *http.Request) error {
	if r.bodyContentType == MultiPartFormDataBody {
		if req.Body == nil {
			return nil
		}

		if err := req.ParseMultipartForm(fileFormLimit); err != nil {
			// empty form
			if err.Error() == "no multipart boundary param in Content-Type" {
				return nil
			}
			return err
		}

		for fieldName, files := range req.MultipartForm.File {
			for _, file := range files {
				f, err := file.Open()
				if err != nil {
					return err
				}
				r.multipartFiles[fieldName] = append(r.multipartFiles[fieldName], File{
					buf:         bufio.NewReader(f),
					contentType: file.Header.Get("Content-Type"),
					name:        file.Filename,
					size:        int(file.Size),
				})
				r.multipartValues[fieldName] = append(r.multipartValues[fieldName], fmt.Sprintf(
					"%s|%d", fieldName, len(r.multipartValues[fieldName])-1,
				))
			}
		}

		for fieldName, values := range req.MultipartForm.Value {
			if _, ok := r.multipartValues[fieldName]; !ok {
				continue
			}

			if len(fieldName) > 0 && len(values) > 0 {
				r.multipartValues[fieldName] = values
			}
		}
	}

	return nil
}

func (r *request) parseUrlEncodedForm(req *http.Request) error {
	if r.bodyContentType == UrlEncodedFormBody {
		if req.Body == nil {
			return nil
		}

		if err := req.ParseForm(); err != nil {
			return err
		}

		r.urlEncodedValues = req.Form
		for fieldName := range req.Form {
			if _, exists := r.query[fieldName]; !exists {
				r.urlEncodedValues.Del(fieldName)
			}

			if len(fieldName) == 0 {
				r.urlEncodedValues.Del(fieldName)
			}
		}
	}

	return nil
}
