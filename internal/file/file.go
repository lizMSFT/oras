/*
Copyright The ORAS Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package file

import (
	"bytes"
	"fmt"
	"io"
	"os"

	digest "github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// PrepareContent prepares the content descriptor from the file path or stdin.
func PrepareContent(path string, mediaType string) (ocispec.Descriptor, io.ReadCloser, error) {
	if path == "" {
		return ocispec.Descriptor{}, nil, fmt.Errorf("missing file name")
	}

	if path == "-" {
		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return ocispec.Descriptor{}, nil, err
		}
		rc := io.NopCloser(bytes.NewReader(content))

		return ocispec.Descriptor{
			MediaType: mediaType,
			Digest:    digest.FromBytes(content),
			Size:      int64(len(content)),
		}, rc, nil
	}

	fp, err := os.Open(path)
	if err != nil {
		return ocispec.Descriptor{}, nil, fmt.Errorf("failed to open %s: %w", path, err)
	}

	fi, err := os.Stat(path)
	if err != nil {
		return ocispec.Descriptor{}, nil, fmt.Errorf("failed to stat %s: %w", path, err)
	}

	dgst, err := digest.FromReader(fp)
	if err != nil {
		return ocispec.Descriptor{}, nil, err
	}

	if _, err = fp.Seek(0, io.SeekStart); err != nil {
		return ocispec.Descriptor{}, nil, err
	}

	return ocispec.Descriptor{
		MediaType: mediaType,
		Digest:    dgst,
		Size:      fi.Size(),
	}, fp, nil
}
