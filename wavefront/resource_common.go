package wavefront_plugin

import (
	"fmt"
	"github.com/WavefrontHQ/go-wavefront-management-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strings"
)

func suppressCase(k, old, new string, d *schema.ResourceData) bool {
	if strings.ToLower(old) == strings.ToLower(new) {
		return true
	}
	return false
}

func suppressSpaces(k, old, new string, d *schema.ResourceData) bool {
	if strings.TrimSpace(old) == strings.TrimSpace(new) {
		return true
	}
	return false
}

func trimSpaces(d interface{}) string {
	if s, ok := d.(string); ok {
		return strings.TrimSpace(s)
	}

	return ""
}

func trimSpacesMap(m map[string]interface{}) map[string]string {
	trimmed := map[string]string{}
	for key, v := range m {
		trimmed[key] = trimSpaces(v)
	}
	return trimmed
}

// Decodes the ACL from the state
func decodeAccessControlList(d *schema.ResourceData) (canView, canModify []string) {
	for _, cv := range d.Get("can_view").(*schema.Set).List() {
		canView = append(canView, cv.(string))
	}

	for _, cv := range d.Get("can_modify").(*schema.Set).List() {
		canModify = append(canModify, cv.(string))
	}

	return canView, canModify
}

// Decodes the tags from the state and returns a []string of tags
func decodeTags(d *schema.ResourceData) (tags []string) {
	for _, tag := range d.Get("tags").(*schema.Set).List() {
		tags = append(tags, tag.(string))
	}

	return tags
}

// Given a GroupID will check if this group is the Everyone group
func isEveryoneGroup(id string, m interface{}) (bool, error) {
	client := m.(*wavefrontClient).client.UserGroups()
	ug := &wavefront.UserGroup{ID: &id}
	err := client.Get(ug)
	if err != nil {
		return false, fmt.Errorf("id provided does not match any user groups. %s", id)
	}

	if ug.Name != "Everyone" {
		return false, nil
	}

	return true, nil
}
