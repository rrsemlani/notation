// Copyright The Notary Project Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package blob

import (
	"errors"
	"fmt"
	"github.com/notaryproject/notation/cmd/notation/internal/osutil"
	"github.com/notaryproject/notation/cmd/notation/internal/outputs"
	"github.com/notaryproject/notation/internal/cmd"
	"github.com/notaryproject/notation/internal/envelope"
	"github.com/spf13/cobra"
	"path/filepath"
)

type blobInspectOpts struct {
	cmd.LoggingFlagOpts
	signaturePath string
	outputFormat  string
}

func inspectCommand(opts *blobInspectOpts) *cobra.Command {
	if opts == nil {
		opts = &blobInspectOpts{}
	}
	longMessage := `Inspect signature associated with the signed blob.

Example - Inspect BLOB signature:
  notation blob inspect <signature_path>

Example - Inspect BLOB signature and output as JSON:
  notation blob inspect --output json <signature_path>
`

	command := &cobra.Command{
		Use:   "blob inspect [signaturePath]",
		Short: "Inspect signature associated with the signed BLOB",
		Long:  longMessage,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("missing signature path: use `notation blob inspect --help` to see what parameters are required")
			}
			opts.signaturePath = args[0]
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBlobInspect(opts)
		},
	}

	opts.LoggingFlagOpts.ApplyFlags(command.Flags())
	cmd.SetPflagOutput(command.Flags(), &opts.outputFormat, cmd.PflagOutputUsage)
	return command
}

func runBlobInspect(opts *blobInspectOpts) error {
	if opts.outputFormat != cmd.OutputJSON && opts.outputFormat != cmd.OutputPlaintext {
		return fmt.Errorf("unrecognized output format %s", opts.outputFormat)
	}

	// initialize
	mediaType, err := envelope.GetEnvelopeMediaType(filepath.Ext(opts.signaturePath))
	if err != nil {
		return err
	}
	signatureEnv, err := osutil.ReadFile(opts.signaturePath, 10485760) // 10Mb in bytes
	if err != nil {
		return err
	}
	output := outputs.InspectOutput{MediaType: mediaType, Signatures: []outputs.SignatureOutput{}}
	err, output.Signatures = outputs.Signatures(mediaType, "", output, signatureEnv)
	if err != nil {
		return nil
	}
	if err := outputs.PrintOutput(opts.outputFormat, opts.signaturePath, output); err != nil {
		return err
	}
	return nil
}
