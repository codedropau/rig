package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	files := []string{
		"./fixtures/docker-compose.yml",
		"./fixtures/docker-compose.osx.yml",
	}
	config, err := Load(files)
	assert.Nil(t, err)

	expected := Config{
		Services: map[string]*Service{
			"nginx": {
				Image: "skpr/nginx:1.x-dev",
				Volumes: []string{
					"data:/data",
					"files:/data/app/sites/default/files",
				},
				Ports: []string{
					"8080:8080",
					"3306:3306",
				},
			},
			"php-fpm": {
				Image: "skpr/php-fpm:7.3-1.x-dev",
				Volumes: []string{
					"data:/data",
					"files:/data/app/sites/default/files",
				},
			},
			"php-fpm-xdebug": {
				Image: "skpr/php-fpm:7.3-1.x-xdebug",
				Volumes: []string{
					"data:/data",
					"files:/data/app/sites/default/files",
				},
				Environment: []string{
					"PHP_FPM_PORT=9001",
					"PHP_IDE_CONFIG=serverName=localhost",
					"XDEBUG_CONFIG=remote_host=host.docker.internal",
				},
			},
			"php-cli": {
				Image: "skpr/php-cli:7.3-1.x-dev",
				Volumes: []string{
					"data:/data",
					"files:/data/app/sites/default/files",
				},
			},
			"php-cli-xdebug": {
				Image: "skpr/php-cli:7.3-1.x-xdebug",
				Volumes: []string{
					"data:/data",
					"files:/data/app/sites/default/files",
				},
				Environment: []string{
					"PHP_IDE_CONFIG=serverName=localhost",
					"XDEBUG_CONFIG=remote_host=host.docker.internal",
				},
			},
			"mysql-default": {
				Image: "646598420362.dkr.ecr.ap-southeast-2.amazonaws.com/skpr-pnx-pnx-d8/mysql:dev-default-latest",
			},
			"mailhog": {
				Image: "mailhog/mailhog",
				Ports: []string{
					"8025:8025",
				},
			},
		},
		Volumes: map[string]*Volume{
			"data": {
				Driver: "local",
				Opts: VolumeOpts{
					Device: "${PWD}",
				},
			},
			"files": {
				Driver: "local",
				Opts: VolumeOpts{
					Device: "${PWD}/app/sites/default/files",
				},
			},
		},
	}

	assert.Equal(t, expected, *config)
}