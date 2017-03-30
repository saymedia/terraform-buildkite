package buildkite

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

const pipelineBadgeURL = "badge_url"
const pipelineBranchConfiguration = "branch_configuration"
const pipelineBuildsURL = "builds_url"
const pipelineCreatedAt = "created_at"
const pipelineDefaultBranch = "default_branch"
const pipelineDescription = "pipelineDescription"
const pipelineEnv = "env"
const pipelineID = "id"
const pipelineName = "name"
const pipelineProviderSettings = "provider_settings"
const pipelineRepository = "repository"
const pipelineSlug = "slug"
const pipelineSteps = "step"
const pipelineURL = "url"
const pipelineWebURL = "web_url"
const pipelineWebhookURL = "webhook_url"
const stepAgentQueryRules = "agent_query_rules"
const stepArtifactPaths = "artifact_paths"
const stepBranchConfiguration = "branch_configuration"
const stepConcurrency = "concurrency"
const stepCommand = "command"
const stepEnv = "env"
const stepName = "name"
const stepParallelism = "parallelism"
const stepTimeoutInMinutes = "timeout_in_minutes"
const stepType = "type"

func resourcePipeline() *schema.Resource {
	return &schema.Resource{
		Create: createPipeline,
		Read:   readPipeline,
		Update: updatePipeline,
		Delete: deletePipeline,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			pipelineID: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			pipelineSlug: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			pipelineWebURL: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			pipelineBuildsURL: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			pipelineCreatedAt: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			pipelineURL: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			pipelineBadgeURL: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			pipelineName: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			pipelineDescription: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			pipelineRepository: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			pipelineBranchConfiguration: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			pipelineDefaultBranch: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			pipelineEnv: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			pipelineProviderSettings: &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeBool,
				},
			},
			pipelineWebhookURL: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			pipelineSteps: &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						stepType: &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						stepName: &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						stepCommand: &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						stepEnv: &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						stepTimeoutInMinutes: &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						stepAgentQueryRules: &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						stepArtifactPaths: &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						stepBranchConfiguration: &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						stepConcurrency: &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						stepParallelism: &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

type pipeline struct {
	BadgeURL            string            `json:"badge_url,omitempty"`
	BranchConfiguration string            `json:"branch_configuration,omitempty"`
	BuildsURL           string            `json:"builds_url,omitempty"`
	CreatedAt           string            `json:"created_at,omitempty"`
	DefaultBranch       string            `json:"default_branch,omitempty"`
	Description         string            `json:"description,omitempty"`
	Environment         map[string]string `json:"env,omitempty"`
	ID                  string            `json:"id,omitempty"`
	Name                string            `json:"name,omitempty"`
	Provider            buildkiteProvider `json:"provider,omitempty"`
	ProviderSettings    map[string]bool   `json:"provider_settings,omitempty"`
	Repository          string            `json:"repository,omitempty"`
	Slug                string            `json:"slug,omitempty"`
	Steps               []step            `json:"steps"`
	URL                 string            `json:"url,omitempty"`
	WebURL              string            `json:"web_url,omitempty"`
}

type buildkiteProvider struct {
	ID         string                 `json:"id"`
	Settings   map[string]interface{} `json:"settings"`
	WebhookURL string                 `json:"webhook_url"`
}

type step struct {
	AgentQueryRules     []string          `json:"agent_query_rules,omitempty"`
	ArtifactPaths       string            `json:"artifact_paths,omitempty"`
	BranchConfiguration string            `json:"branch_configuration,omitempty"`
	Command             string            `json:"command,omitempty"`
	Concurrency         int               `json:"concurrency,omitempty"`
	Environment         map[string]string `json:"env,omitempty"`
	Name                string            `json:"name,omitempty"`
	Parallelism         int               `json:"parallelism,omitempty"`
	TimeoutInMinutes    int               `json:"timeout_in_minutes,omitempty"`
	Type                string            `json:"type"`
}

func createPipeline(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] CreatePipeline")

	client := meta.(*Client)

	req := preparePipelineRequestPayload(d)
	res := &pipeline{}

	err := client.Post([]string{"pipelines"}, req, res)
	if err != nil {
		return err
	}

	updatePipelineFromAPI(d, res)

	return nil
}

func readPipeline(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] ReadPipeline")

	client := meta.(*Client)
	slug := d.Id()

	res := &pipeline{}

	err := client.Get([]string{"pipelines", slug}, res)
	if err != nil {
		if _, ok := err.(*notFound); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	updatePipelineFromAPI(d, res)

	return nil
}

func updatePipeline(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] UpdatePipeline")

	client := meta.(*Client)
	slug := d.Id()

	req := preparePipelineRequestPayload(d)
	res := &pipeline{}

	err := client.Patch([]string{"pipelines", slug}, req, res)
	if err != nil {
		return err
	}

	updatePipelineFromAPI(d, res)

	return nil
}

func deletePipeline(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[TRACE] DeletePipeline")

	client := meta.(*Client)

	slug := d.Id()

	return client.Delete([]string{"pipelines", slug})
}

func updatePipelineFromAPI(d *schema.ResourceData, p *pipeline) {
	d.SetId(p.Slug)

	d.Set(pipelineBuildsURL, p.BuildsURL)
	d.Set(pipelineBranchConfiguration, p.BranchConfiguration)
	d.Set(pipelineDefaultBranch, p.DefaultBranch)
	d.Set(pipelineDescription, p.Description)
	d.Set(pipelineEnv, p.Environment)
	d.Set(pipelineID, p.ID)
	d.Set(pipelineName, p.Name)
	d.Set(pipelineProviderSettings, p.Provider.Settings)
	d.Set(pipelineRepository, p.Repository)
	d.Set(pipelineSlug, p.Slug)
	d.Set(pipelineWebhookURL, p.Provider.WebhookURL)
	d.Set(pipelineWebURL, p.WebURL)

	stepMap := make([]interface{}, len(p.Steps))
	for i, vI := range p.Steps {
		stepMap[i] = map[string]interface{}{
			stepAgentQueryRules:     vI.AgentQueryRules,
			stepArtifactPaths:       vI.ArtifactPaths,
			stepBranchConfiguration: vI.BranchConfiguration,
			stepCommand:             vI.Command,
			stepConcurrency:         vI.Concurrency,
			stepEnv:                 vI.Environment,
			stepName:                vI.Name,
			stepParallelism:         vI.Parallelism,
			stepTimeoutInMinutes:    vI.TimeoutInMinutes,
			stepType:                vI.Type,
		}
	}
	d.Set(pipelineSteps, stepMap)
}

func preparePipelineRequestPayload(d *schema.ResourceData) *pipeline {
	req := &pipeline{}

	req.BranchConfiguration = d.Get(pipelineBranchConfiguration).(string)
	req.DefaultBranch = d.Get(pipelineDefaultBranch).(string)
	req.Description = d.Get(pipelineDescription).(string)
	req.Environment = map[string]string{}
	for k, vI := range d.Get(pipelineEnv).(map[string]interface{}) {
		req.Environment[k] = vI.(string)
	}
	req.Name = d.Get(pipelineName).(string)
	req.ProviderSettings = map[string]bool{}
	for k, vI := range d.Get(pipelineProviderSettings).(map[string]interface{}) {
		req.ProviderSettings[k] = vI.(bool)
	}
	req.Repository = d.Get(pipelineRepository).(string)
	req.Slug = d.Get(pipelineSlug).(string)

	stepsI := d.Get(pipelineSteps).([]interface{})
	req.Steps = make([]step, len(stepsI))

	for i, stepI := range stepsI {
		stepM := stepI.(map[string]interface{})
		req.Steps[i] = step{
			AgentQueryRules:     make([]string, len(stepM[stepAgentQueryRules].([]interface{}))),
			ArtifactPaths:       stepM[stepArtifactPaths].(string),
			BranchConfiguration: stepM[stepBranchConfiguration].(string),
			Command:             stepM[stepCommand].(string),
			Concurrency:         stepM[stepConcurrency].(int),
			Environment:         map[string]string{},
			Name:                stepM[stepName].(string),
			Parallelism:         stepM[stepParallelism].(int),
			TimeoutInMinutes:    stepM[stepTimeoutInMinutes].(int),
			Type:                stepM[stepType].(string),
		}

		for j, vI := range stepM[stepAgentQueryRules].([]interface{}) {
			req.Steps[i].AgentQueryRules[j] = vI.(string)
		}

		for k, vI := range stepM[stepEnv].(map[string]interface{}) {
			req.Steps[i].Environment[k] = vI.(string)
		}
	}

	return req
}
