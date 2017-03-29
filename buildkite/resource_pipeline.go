package buildkite

import (
    "github.com/hashicorp/terraform/helper/schema"
)

func resourcePipeline() *schema.Resource {
    return &schema.Resource{
        Create: CreatePipeline,
        Read:   ReadPipeline,
        Update: UpdatePipeline,
        Delete: DeletePipeline,
        Importer: &schema.ResourceImporter{
            State: schema.ImportStatePassthrough,
        },

        Schema: map[string]*schema.Schema{
            "id": &schema.Schema{
                Type:     schema.TypeString,
                Computed: true,
            },
            "slug": &schema.Schema{
                Type:     schema.TypeString,
                Computed: true,
                Optional: true,
            },
            "web_url": &schema.Schema{
                Type:     schema.TypeString,
                Computed: true,
            },
            "builds_url": &schema.Schema{
                Type:     schema.TypeString,
                Computed: true,
            },
            "created_at": &schema.Schema{
                Type:     schema.TypeString,
                Computed: true,
            },
            "url": &schema.Schema{
                Type:     schema.TypeString,
                Computed: true,
            },
            "badge_url": &schema.Schema{
                Type:     schema.TypeString,
                Computed: true,
            },
            "name": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
            },
            "description": &schema.Schema{
                Type:     schema.TypeString,
                Optional: true,
            },
            "repository": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
            },
            "branch_configuration": &schema.Schema{
                Type:     schema.TypeString,
                Optional: true,
            },
            "env": &schema.Schema{
                Type:     schema.TypeMap,
                Optional: true,
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "provider_settings": &schema.Schema{
                Type:     schema.TypeMap,
                Optional: true,
                Elem: &schema.Schema{
                    Type: schema.TypeBool,
                },
            },
            "webhook_url": &schema.Schema{
                Type:     schema.TypeString,
                Computed: true,
            },
            "step": &schema.Schema{
                Type:     schema.TypeList,
                Required: true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "type": &schema.Schema{
                            Type:     schema.TypeString,
                            Required: true,
                        },
                        "name": &schema.Schema{
                            Type:     schema.TypeString,
                            Optional: true,
                        },
                        "command": &schema.Schema{
                            Type:     schema.TypeString,
                            Optional: true,
                        },
                        "env": &schema.Schema{
                            Type:     schema.TypeMap,
                            Optional: true,
                            Elem: &schema.Schema{
                                Type: schema.TypeString,
                            },
                        },
                        "timeout_in_minutes": &schema.Schema{
                            Type:     schema.TypeInt,
                            Optional: true,
                        },
                        "agent_query_rules": &schema.Schema{
                            Type:     schema.TypeList,
                            Optional: true,
                            Elem: &schema.Schema{
                                Type: schema.TypeString,
                            },
                        },
                        "artifact_paths": &schema.Schema{
                            Type:     schema.TypeString,
                            Optional: true,
                        },
                        "branch_configuration": &schema.Schema{
                            Type:     schema.TypeString,
                            Optional: true,
                        },
                        "concurrency": &schema.Schema{
                            Type:     schema.TypeInt,
                            Optional: true,
                        },
                        "parallelism": &schema.Schema{
                            Type:     schema.TypeInt,
                            Optional: true,
                        },
                    },
                },
            },
        },
    }
}

type Pipeline struct {
    Id                  string            `json:"id,omitempty"`
    Environment         map[string]string `json:"env,omitempty"`
    Slug                string            `json:"slug,omitempty"`
    WebURL              string            `json:"web_url,omitempty"`
    BuildsURL           string            `json:"builds_url,omitempty"`
    Url                 string            `json:"url,omitempty"`
    BadgeURL            string            `json:"badge_url,omitempty"`
    CreatedAt           string            `json:"created_at,omitempty"`
    Repository          string            `json:"repository,omitempty"`
    Name                string            `json:"name,omitempty"`
    Description         string            `json:"description,omitempty"`
    BranchConfiguration string            `json:"branch_configuration,omitempty"`
    Provider            BuildkiteProvider `json:"provider,omitempty"`
    ProviderSettings    map[string]bool   `json:"provider_settings,omitempty"`
    Steps               []Step            `json:"steps"`
}

type BuildkiteProvider struct {
    Id         string                 `json:"id"`
    Settings   map[string]interface{} `json:"settings"`
    WebhookURL string                 `json:"webhook_url"`
}

type Step struct {
    Type                string            `json:"type"`
    Name                string            `json:"name,omitempty"`
    Command             string            `json:"command,omitempty"`
    Environment         map[string]string `json:"env,omitempty"`
    TimeoutInMinutes    int               `json:"timeout_in_minutes,omitempty"`
    AgentQueryRules     []string          `json:"agent_query_rules,omitempty"`
    BranchConfiguration string            `json:"branch_configuration,omitempty"`
    ArtifactPaths       string            `json:"artifact_paths,omitempty"`
    Concurrency         int               `json:"concurrency,omitempty"`
    Parallelism         int               `json:"parallelism,omitempty"`
}

func CreatePipeline(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*Client)

    req := preparePipelineRequestPayload(d)
    res := &Pipeline{}

    err := client.Post([]string{"pipelines"}, req, res)
    if err != nil {
        return err
    }

    updatePipelineFromAPI(d, res)

    return nil
}

func ReadPipeline(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*Client)
    slug := d.Id()

    res := &Pipeline{}

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

func UpdatePipeline(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*Client)
    slug := d.Id()

    req := preparePipelineRequestPayload(d)
    res := &Pipeline{}

    err := client.Patch([]string{"pipelines", slug}, req, res)
    if err != nil {
        return err
    }

    updatePipelineFromAPI(d, res)

    return nil
}

func DeletePipeline(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*Client)

    slug := d.Id()

    return client.Delete([]string{"pipelines", slug})
}

func updatePipelineFromAPI(d *schema.ResourceData, p *Pipeline) {
    d.SetId(p.Slug)
    d.Set("id", p.Id)
    d.Set("env", p.Environment)
    d.Set("name", p.Name)
    d.Set("description", p.Description)
    d.Set("repository", p.Repository)
    d.Set("web_url", p.WebURL)
    d.Set("slug", p.Slug)
    d.Set("builds_url", p.BuildsURL)
    d.Set("branch_configuration", p.BranchConfiguration)
    d.Set("provider_settings", p.Provider.Settings)
    d.Set("webhook_url", p.Provider.WebhookURL)
}

func preparePipelineRequestPayload(d *schema.ResourceData) *Pipeline {
    req := &Pipeline{}

    req.Name = d.Get("name").(string)
    req.Description = d.Get("description").(string)
    req.Slug = d.Get("slug").(string)
    req.Repository = d.Get("repository").(string)
    req.BranchConfiguration = d.Get("branch_configuration").(string)
    req.Environment = map[string]string{}
    for k, vI := range d.Get("env").(map[string]interface{}) {
        req.Environment[k] = vI.(string)
    }
    req.ProviderSettings = map[string]bool{}
    for k, vI := range d.Get("provider_settings").(map[string]interface{}) {
        req.ProviderSettings[k] = vI.(bool)
    }

    stepsI := d.Get("step").([]interface{})
    req.Steps = make([]Step, len(stepsI))

    for i, stepI := range stepsI {
        stepM := stepI.(map[string]interface{})
        req.Steps[i] = Step{
            Type:                stepM["type"].(string),
            Name:                stepM["name"].(string),
            Command:             stepM["command"].(string),
            Environment:         map[string]string{},
            AgentQueryRules:     make([]string, len(stepM["agent_query_rules"].([]interface{}))),
            BranchConfiguration: stepM["branch_configuration"].(string),
            ArtifactPaths:       stepM["artifact_paths"].(string),
            Concurrency:         stepM["concurrency"].(int),
            Parallelism:         stepM["parallelism"].(int),
        }

        for k, vI := range stepM["env"].(map[string]interface{}) {
            req.Steps[i].Environment[k] = vI.(string)
        }

        for j, vI := range stepM["agent_query_rules"].([]interface{}) {
            req.Steps[i].AgentQueryRules[j] = vI.(string)
        }
    }

    return req
}
