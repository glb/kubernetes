/*
Copyright 2016 The Kubernetes Authors.

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

package dockershim

import (
	"testing"

	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"

	runtimeapi "k8s.io/kubernetes/pkg/kubelet/apis/cri/runtime/v1alpha2"
)

func TestConvertDockerStatusToRuntimeAPIState(t *testing.T) {
	testCases := []struct {
		input    string
		expected runtimeapi.ContainerState
	}{
		{input: "Up 5 hours", expected: runtimeapi.ContainerState_CONTAINER_RUNNING},
		{input: "Exited (0) 2 hours ago", expected: runtimeapi.ContainerState_CONTAINER_EXITED},
		{input: "Created", expected: runtimeapi.ContainerState_CONTAINER_CREATED},
		{input: "Random string", expected: runtimeapi.ContainerState_CONTAINER_UNKNOWN},
	}

	for _, test := range testCases {
		actual := toRuntimeAPIContainerState(test.input)
		assert.Equal(t, test.expected, actual)
	}
}

func TestConvertToPullableImageID(t *testing.T) {
	testCases := []struct {
		id       string
		image    *dockertypes.ImageInspect
		expected string
	}{
		{
			id: "image-1",
			image: &dockertypes.ImageInspect{
				RepoDigests: []string{"digest-1"},
			},
			expected: DockerPullableImageIDPrefix + "digest-1",
		},
		{
			id: "image-2",
			image: &dockertypes.ImageInspect{
				RepoDigests: []string{},
			},
			expected: DockerImageIDPrefix + "image-2",
		},
	}

	for _, test := range testCases {
		actual := toPullableImageID(test.id, test.image)
		assert.Equal(t, test.expected, actual)
	}
}

func TestImageInspectToRuntimeAPIImage_ImageLabels(t *testing.T) {
	testCases := []struct {
		id       string
		image    *dockertypes.ImageInspect
		expected map[string]string
	}{
		{
			id: "image labels are copied from the image config to the runtime API image",
			image: &dockertypes.ImageInspect{
				Config: &container.Config{
					Labels: map[string]string{
						"key": "value",
					},
				},
			},
			expected: map[string]string{
				"key": "value",
			},
		},
	}

	for _, test := range testCases {
		actual, err := imageInspectToRuntimeAPIImage(test.image)
		if err != nil {
			t.Errorf("[%s] unexpected error %T %v", test.id, err, err)
			continue
		}
		assert.Equal(t, test.expected, actual.ImageLabels)
	}
}

func TestImageToRuntimeAPIImage_ImageLabels(t *testing.T) {
	testCases := []struct {
		id       string
		image    *dockertypes.ImageSummary
		expected map[string]string
	}{
		{
			id: "image labels are copied from the image summary to the runtime API image",
			image: &dockertypes.ImageSummary{
				Labels: map[string]string{
					"key": "value",
				},
			},
			expected: map[string]string{
				"key": "value",
			},
		},
	}

	for _, test := range testCases {
		actual, err := imageToRuntimeAPIImage(test.image)
		if err != nil {
			t.Errorf("[%s] unexpected error %T %v", test.id, err, err)
			continue
		}
		assert.Equal(t, test.expected, actual.ImageLabels)
	}
}
