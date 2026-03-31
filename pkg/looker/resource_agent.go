package looker

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apiclient "github.com/looker-open-source/sdk-codegen/go/sdk/v4"
)

func resourceAgent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAgentCreate,
		ReadContext:   resourceAgentRead,
		UpdateContext: resourceAgentUpdate,
		DeleteContext: resourceAgentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Agent name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Agent description",
			},
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The category of the agent (e.g., dashboard, conversation)",
			},
			"sources": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Agent sources",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"model": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Source model",
						},
						"explore": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Source explore",
						},
					},
				},
			},
			"instructions": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Agent instructions (context)",
			},
			"code_interpreter": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enables Code Interpreter for this Agent",
			},
		},
	}
}

func resourceAgentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	agentName := d.Get("name").(string)

	writeAgent := buildWriteAgent(d)

	agent, err := client.CreateAgent(writeAgent, "", nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "CreateAgent", "agent", "%s", agentName))
	}

	if agent.Id == nil {
		return diag.Errorf("Agent ID not returned from API")
	}

	d.SetId(*agent.Id)

	return resourceAgentRead(ctx, d, m)
}

func resourceAgentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	agentID := d.Id()

	agent, err := client.GetAgent(agentID, "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(wrapSDKError(err, "GetAgent", "agent", "%s", agentID))
	}

	if agent.Name != nil {
		if err := d.Set("name", *agent.Name); err != nil {
			return diag.FromErr(err)
		}
	}

	if agent.Description != nil {
		if err := d.Set("description", *agent.Description); err != nil {
			return diag.FromErr(err)
		}
	}

	if agent.Category != nil {
		if err := d.Set("category", *agent.Category); err != nil {
			return diag.FromErr(err)
		}
	}

	if agent.Sources != nil {
		sources := make([]map[string]interface{}, len(*agent.Sources))
		for i, s := range *agent.Sources {
			source := map[string]interface{}{}
			if s.Model != nil {
				source["model"] = *s.Model
			}
			if s.Explore != nil {
				source["explore"] = *s.Explore
			}
			sources[i] = source
		}
		if err := d.Set("sources", sources); err != nil {
			return diag.FromErr(err)
		}
	}

	if agent.Context != nil && agent.Context.Instructions != nil {
		if err := d.Set("instructions", *agent.Context.Instructions); err != nil {
			return diag.FromErr(err)
		}
	}

	if agent.CodeInterpreter != nil {
		if err := d.Set("code_interpreter", *agent.CodeInterpreter); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceAgentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	agentID := d.Id()

	writeAgent := buildWriteAgent(d)

	_, err := client.UpdateAgent(agentID, writeAgent, "", nil)
	if err != nil {
		return diag.FromErr(wrapSDKError(err, "UpdateAgent", "agent", "id=%s", agentID))
	}

	return resourceAgentRead(ctx, d, m)
}

func resourceAgentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*apiclient.LookerSDK)

	agentID := d.Id()
	agentName := d.Get("name").(string)

	_, err := client.DeleteAgent(agentID, "", nil)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil
		}
		return diag.FromErr(wrapSDKError(err, "DeleteAgent", "agent", "name=%s, id=%s", agentName, agentID))
	}

	return nil
}

func buildWriteAgent(d *schema.ResourceData) apiclient.WriteAgent {
	agentName := d.Get("name").(string)
	description := d.Get("description").(string)
	codeInterpreter := d.Get("code_interpreter").(bool)

	writeAgent := apiclient.WriteAgent{
		Name:            &agentName,
		Description:     &description,
		CodeInterpreter: &codeInterpreter,
	}

	if v, ok := d.GetOk("category"); ok {
		category := v.(string)
		writeAgent.Category = &category
	}

	if v, ok := d.GetOk("sources"); ok {
		sourcesList := v.([]interface{})
		sources := make([]apiclient.Source, len(sourcesList))
		for i, raw := range sourcesList {
			s := raw.(map[string]interface{})
			model := s["model"].(string)
			explore := s["explore"].(string)
			sources[i] = apiclient.Source{
				Model:   &model,
				Explore: &explore,
			}
		}
		writeAgent.Sources = &sources
	}

	if v, ok := d.GetOk("instructions"); ok {
		instructions := v.(string)
		writeAgent.Context = &apiclient.Context{
			Instructions: &instructions,
		}
	}

	return writeAgent
}
