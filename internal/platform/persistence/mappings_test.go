package persistence

import (
	"context"
	"github.com/b2wdigital/restQL-golang/internal/domain"
	"github.com/b2wdigital/restQL-golang/internal/platform/logger"
	"github.com/b2wdigital/restQL-golang/test"
	"io/ioutil"
	"testing"
)

const defaultTenant = "default"

func TestMappingsReader_Env(t *testing.T) {
	envSource := stubEnvSource{
		getAll: map[string]string{
			"RESTQL_MAPPING_HERO":     "http://hero.api/",
			"RESTQL_MAPPING_SIDEKICK": "http://sidekick.api/",
			"TEST_VAR":                "foo",
		},
	}
	db := stubDatabase{}

	reader := NewMappingReader(noOpLogger, envSource, map[string]string{}, db)

	expected := map[string]domain.Mapping{
		"hero": {
			ResourceName:  "hero",
			Schema:        "http",
			Host:          "hero.api",
			Path:          "/",
			PathParams:    []string{},
			PathParamsSet: map[string]struct{}{},
		},
		"sidekick": {
			ResourceName:  "sidekick",
			Schema:        "http",
			Host:          "sidekick.api",
			Path:          "/",
			PathParams:    []string{},
			PathParamsSet: map[string]struct{}{},
		},
	}

	mappings, err := reader.FromTenant(context.Background(), defaultTenant)

	test.VerifyError(t, err)
	test.Equal(t, mappings, expected)
}

func TestMappingsReader_Local(t *testing.T) {
	envSource := stubEnvSource{getAll: map[string]string{}}
	local := map[string]string{
		"hero":     "http://hero.api/",
		"sidekick": "http://sidekick.api/",
	}
	db := stubDatabase{}

	reader := NewMappingReader(noOpLogger, envSource, local, db)

	expected := map[string]domain.Mapping{
		"hero": {
			ResourceName:  "hero",
			Schema:        "http",
			Host:          "hero.api",
			Path:          "/",
			PathParams:    []string{},
			PathParamsSet: map[string]struct{}{},
		},
		"sidekick": {
			ResourceName:  "sidekick",
			Schema:        "http",
			Host:          "sidekick.api",
			Path:          "/",
			PathParams:    []string{},
			PathParamsSet: map[string]struct{}{},
		},
	}

	mappings, err := reader.FromTenant(context.Background(), defaultTenant)

	test.VerifyError(t, err)
	test.Equal(t, mappings, expected)
}

func TestMappingsReader_Database(t *testing.T) {
	envSource := stubEnvSource{getAll: map[string]string{}}
	local := map[string]string{}
	db := stubDatabase{findMappingsForTenant: []domain.Mapping{
		{
			ResourceName:  "hero",
			Schema:        "http",
			Host:          "hero.api",
			Path:          "/",
			PathParams:    []string{},
			PathParamsSet: map[string]struct{}{},
		},
		{
			ResourceName:  "sidekick",
			Schema:        "http",
			Host:          "sidekick.api",
			Path:          "/",
			PathParams:    []string{},
			PathParamsSet: map[string]struct{}{},
		},
	}}

	reader := NewMappingReader(noOpLogger, envSource, local, db)

	expected := map[string]domain.Mapping{
		"hero": {
			ResourceName:  "hero",
			Schema:        "http",
			Host:          "hero.api",
			Path:          "/",
			PathParams:    []string{},
			PathParamsSet: map[string]struct{}{},
		},
		"sidekick": {
			ResourceName:  "sidekick",
			Schema:        "http",
			Host:          "sidekick.api",
			Path:          "/",
			PathParams:    []string{},
			PathParamsSet: map[string]struct{}{},
		},
	}

	mappings, err := reader.FromTenant(context.Background(), defaultTenant)

	test.VerifyError(t, err)
	test.Equal(t, mappings, expected)
}

func TestMappingsReader_ShouldOverwriteMappings(t *testing.T) {
	local := map[string]string{
		"hero":     "http://hero.api/",
		"sidekick": "http://sidekick.api/",
		"villain":  "http://villain.api/",
	}
	db := stubDatabase{
		findMappingsForTenant: []domain.Mapping{
			{
				ResourceName:  "sidekick",
				Schema:        "https",
				Host:          "sidekick.com",
				Path:          "/api",
				PathParams:    []string{},
				PathParamsSet: map[string]struct{}{},
			},
		},
	}
	envSource := stubEnvSource{
		getAll: map[string]string{
			"RESTQL_MAPPING_HERO": "https://hero.com/api/",
		},
	}

	reader := NewMappingReader(noOpLogger, envSource, local, db)

	expected := map[string]domain.Mapping{
		"hero": {
			ResourceName:  "hero",
			Schema:        "https",
			Host:          "hero.com",
			Path:          "/api/",
			PathParams:    []string{},
			PathParamsSet: map[string]struct{}{},
		},
		"sidekick": {
			ResourceName:  "sidekick",
			Schema:        "https",
			Host:          "sidekick.com",
			Path:          "/api",
			PathParams:    []string{},
			PathParamsSet: map[string]struct{}{},
		},
		"villain": {
			ResourceName:  "villain",
			Schema:        "http",
			Host:          "villain.api",
			Path:          "/",
			PathParams:    []string{},
			PathParamsSet: map[string]struct{}{},
		},
	}

	mappings, err := reader.FromTenant(context.Background(), defaultTenant)

	test.VerifyError(t, err)
	test.Equal(t, mappings, expected)
}

var noOpLogger = logger.New(ioutil.Discard, logger.LogOptions{})

type stubDatabase struct {
	findMappingsForTenant []domain.Mapping
	findQuery             domain.SavedQuery
}

func (s stubDatabase) FindMappingsForTenant(ctx context.Context, tenantId string) ([]domain.Mapping, error) {
	return s.findMappingsForTenant, nil
}

func (s stubDatabase) FindQuery(ctx context.Context, namespace string, name string, revision int) (domain.SavedQuery, error) {
	return s.findQuery, nil
}

type stubEnvSource struct {
	getAll map[string]string
}

func (s stubEnvSource) GetString(key string) string {
	return ""
}

func (s stubEnvSource) GetAll() map[string]string {
	return s.getAll
}
