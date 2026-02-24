---
id: {{.ID}}
title: {{.Title}}
kind: {{.Kind}}
performed_by: {{.PerformedBy}}
my_role: {{.MyRole}}
work_id: {{.WorkID}}
context:
  label: {{.Context.Label}}
  kind: {{.Context.Kind}}
  location: {{.Context.Location}}
  year: {{.Context.Year}}
date_start: {{.DateStart}}
date_end: {{.DateEnd}}
description: {{.Description}}
collaborators: []
attachments: []
---
{{.ContentBody}}
