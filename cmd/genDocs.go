// Copyright © 2019 Christian Rolly <mail@chromium-solutions.de>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// genDocsCmd represents the genDocs command
var (
	genDocsCmd = &cobra.Command{
		Use:   "genDocs",
		Short: "Generates the documentation for mug",
		Run: func(cmd *cobra.Command, args []string) {
			const fmTemplate = `---
date: %s
title: "%s"
slug: %s
url: %s
---
`

			filePrepender := func(filename string) string {
				now := time.Now().Format(time.RFC3339)
				name := filepath.Base(filename)
				base := strings.TrimSuffix(name, path.Ext(name))
				url := strings.ToLower(base) + "/"
				return fmt.Sprintf(fmTemplate, now, strings.Replace(base, "_", " ", -1), base, url)
			}

			linkHandler := func(name string) string {
				base := strings.TrimSuffix(name, path.Ext(name))
				return strings.ToLower(base) + "/"
			}

			err := doc.GenMarkdownTreeCustom(RootCmd, "./docs", filePrepender, linkHandler)
			if err != nil {
				log.Fatal(err)
			}

		},
	}
)

func init() {
	RootCmd.AddCommand(genDocsCmd)
}
