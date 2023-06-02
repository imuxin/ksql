package envtest

import (
	"context"
	"embed"
	"io/fs"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/yaml"

	"github.com/imuxin/ksql/pkg/repl"
)

func getAllFilenames(efs *embed.FS) (files []string, err error) {
	if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		files = append(files, path)

		return nil
	}); err != nil {
		return nil, err
	}

	return files, nil
}

//go:embed testdata/*
var content embed.FS

type resource struct {
	APIVersion string                 `json:"apiVersion"`
	Kind       string                 `json:"kind"`
	Metadata   map[string]interface{} `json:"metadata"`
}

func prepareDatabase(t *testing.T, restConfig *rest.Config) {
	restConfig.APIPath = "api"
	restConfig.NegotiatedSerializer = scheme.Codecs
	data, _ := content.ReadFile("testdata/data.yaml")
	for _, item := range strings.Split(string(data), "---") {
		r := &resource{}
		assert.NoError(t, utilyaml.Unmarshal([]byte(item), r))
		gv, err := schema.ParseGroupVersion(r.APIVersion)
		assert.NoError(t, err)
		restConfig.GroupVersion = &gv
		restCli, err := rest.RESTClientFor(restConfig)
		assert.NoError(t, err)

		jsonBody, err := yaml.YAMLToJSON([]byte(item))
		assert.NoError(t, err)
		req := restCli.Post().Resource(strings.ToLower(r.Kind) + "s").Body(jsonBody)
		if r.Metadata["namespace"] != nil {
			req = req.Namespace(r.Metadata["namespace"].(string))
		}
		err = req.Do(context.TODO()).Error()
		assert.NoError(t, err)
	}
}

func TestExecuteAndFormat(t *testing.T) {
	env := &envtest.Environment{}
	restConfig, err := env.Start()
	assert.NoError(t, err)

	prepareDatabase(t, restConfig)

	files, err := getAllFilenames(&content)
	assert.NoError(t, err)

	sqls := lo.Filter(files, func(item string, index int) bool {
		return strings.HasSuffix(item, ".sql")
	})
	lo.ForEach(sqls, func(item string, _ int) {
		b, err := content.ReadFile(item)
		assert.NoError(t, err)

		// columns, result, err := executor.ExecuteLikeSQL[ext.Object](string(b), restConfig)
		repl.Exec(string(b), restConfig)
		expect, err := content.ReadFile(strings.TrimSuffix(item, ".sql") + ".output")
		assert.NoError(t, err)
		// fmt.Println(repl.Format(columns, result))
		assert.Equal(t, string(expect), repl.Exec(string(b), restConfig))
	})

	assert.NoError(t, env.Stop())
}
