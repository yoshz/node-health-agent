{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "node-health-agent.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "node-health-agent.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "node-health-agent.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "node-health-agent.labels" -}}
{{ include "node-health-agent.selectorLabels" . }}
chart: {{ include "node-health-agent.chart" . }}
{{- if .Chart.AppVersion }}
version: {{ .Chart.AppVersion | quote }}
{{- end }}
release: {{ .Release.Service }}
managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "node-health-agent.selectorLabels" -}}
app: {{ include "node-health-agent.name" . }}
instance: {{ .Release.Name }}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "node-health-agent.rbacName" -}}
{{- if .Values.rbac.create -}}
    {{ default (include "node-health-agent.fullname" .) .Values.rbac.name }}
{{- else -}}
    {{ default "default" .Values.rbac.name }}
{{- end -}}
{{- end -}}

{{/*
Create the image url
*/}}
{{- define "node-health-agent.image" -}}
{{- .Values.image.repository -}}
{{- if .Values.image.name -}}/{{- .Values.image.name -}}{{- end -}}
{{- .Values.image.tag -}}
{{- end -}}
