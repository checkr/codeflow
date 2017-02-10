package plugins

import "encoding/json"

func ProjectCreateMock() Project {
	var payload Project
	s := `{
		"action": "create",
		"slug": "checkr-deploy-test",
		"repository": "checkr/deploy-test"
	}`
	_ = json.Unmarshal([]byte(s), &payload)
	return payload
}

func GitPingMock() GitPing {
	var payload GitPing
	s := `{
		"repository": "checkr/deploy-test",
		"user": "sasso"
	}`
	_ = json.Unmarshal([]byte(s), &payload)
	return payload
}

func GitCommitMock() GitCommit {
	var payload GitCommit
	s := `{
		"hash": "b4ffef4b5b30b33f024c805cc0f731c5298a4a9f",
		"message": "Update README.md",
		"parentHash": "402cfbbfb4601c01ba91837c46229c2d01f010ae",
		"ref": "refs/heads/master",
		"repository": "checkr/deploy-test",
		"user": "sasso"
	}`
	_ = json.Unmarshal([]byte(s), &payload)
	return payload
}

func GitStatusCirclePendingMock() GitStatus {
	var payload GitStatus
	s := `{
	"context": "ci/circleci",
	"hash": "b4ffef4b5b30b33f024c805cc0f731c5298a4a9f",
	"repository": "checkr/deploy-test",
	"state": "pending",
	"user": "sasso"
	}`
	_ = json.Unmarshal([]byte(s), &payload)
	return payload
}

func GitStatusCircleFailedMock() GitStatus {
	var payload GitStatus
	s := `{
	"context": "ci/circleci",
	"hash": "b4ffef4b5b30b33f024c805cc0f731c5298a4a9f",
	"repository": "checkr/deploy-test",
	"state": "failed",
	"user": "sasso"
	}`
	_ = json.Unmarshal([]byte(s), &payload)
	return payload
}

func GitStatusCircleSuccessMock() GitStatus {
	var payload GitStatus
	s := `{
	"context": "ci/circleci",
	"hash": "b4ffef4b5b30b33f024c805cc0f731c5298a4a9f",
	"repository": "checkr/deploy-test",
	"state": "success",
	"user": "sasso"
	}`
	_ = json.Unmarshal([]byte(s), &payload)
	return payload
}

func DockerDeployCreateMock() DockerDeploy {
	var payload DockerDeploy
	s := `{
		"action": "create",
		"deploymentStrategy": "",
		"docker": {
			"image": "saso/saso",
			"registry": {
				"email": "",
				"host": "",
				"password": "",
				"username": ""
			}
		},
		"environment": "development",
		"project": {
			"repository": "checkr/deploy-test",
			"slug": "checkr-deploy-test",
			"notifyChannels": ["#eng-deploys", "#devops"]
		},
		"release": {
			"headFeature": {
				"hash": "b4ffef4b5b30b33f024c805cc0f731c5298a4a9f",
				"message": "Update README.md",
				"parentHash": "402cfbbfb4601c01ba91837c46229c2d01f010ae",
				"user": "sasso"
			},
			"id": "589df9dacd503701e66f5b2f",
			"tailFeature": {
				"hash": "b4ffef4b5b30b33f024c805cc0f731c5298a4a9f",
				"message": "Update README.md",
				"parentHash": "402cfbbfb4601c01ba91837c46229c2d01f010ae",
				"user": "sasso"
			}
		},
		"secrets": [
			{
				"key": "CODEFLOW_SLUG",
				"type": "env",
				"value": "checkr-deploy-test"
			},
			{
				"key": "CODEFLOW_HASH",
				"type": "env",
				"value": "b4ffef4"
			},
			{
				"key": "CODEFLOW_CREATED_AT",
				"type": "env",
				"value": "2017-02-10T17:35:22Z"
			}
		],
		"services": [
			{
				"action": "create",
				"command": "www",
				"listeners": null,
				"name": "www",
				"replicas": 1,
				"spec": {
					"cpuLimit": "2000m",
					"cpuRequest": "600m",
					"memoryLimit": "2Gi",
					"memoryRequest": "1Gi",
					"terminationGracePeriodSeconds": 300
				},
				"state": "",
				"stateMessage": ""
			}
		],
		"state": "",
		"stateMessage": "",
		"timeout": 0
	}`
	if err := json.Unmarshal([]byte(s), &payload); err != nil {
		panic(err)
	}
	return payload
}

func DockerBuildMock() {

}

func LoadBalancerMock() {

}

func WebsocketMsgMock() {

}
