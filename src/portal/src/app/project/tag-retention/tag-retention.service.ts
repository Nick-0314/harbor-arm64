// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
import { Injectable } from "@angular/core";
import { HttpClient, HttpParams, HttpResponse } from "@angular/common/http";
import { Retention, RuleMetadate } from "./retention";
import { Observable, throwError as observableThrowError } from "rxjs";
import { map, catchError } from "rxjs/operators";
import { Project } from "../project";
import { buildHttpRequestOptionsWithObserveResponse } from "@harbor/ui";

@Injectable()
export class TagRetentionService {
    private I18nMap: object = {
        "retain": "ACTION_RETAIN",
        "lastXDays": "RULE_NAME_1",
        "latestActiveK": "RULE_NAME_2",
        "latestPushedK": "RULE_NAME_3",
        "latestPulledN": "RULE_NAME_4",
        "always": "RULE_NAME_5",
        "nDaysSinceLastPull": "RULE_NAME_6",
        "nDaysSinceLastPush": "RULE_NAME_7",
        "the images from the last # days": "RULE_TEMPLATE_1",
        "the most recent active # images": "RULE_TEMPLATE_2",
        "the most recently pushed # images": "RULE_TEMPLATE_3",
        "the most recently pulled # images": "RULE_TEMPLATE_4",
        "pulled within the last # days": "RULE_TEMPLATE_6",
        "pushed within the last # days": "RULE_TEMPLATE_7",
        "repoMatches": "MAT",
        "repoExcludes": "EXC",
        "matches": "MAT",
        "excludes": "EXC",
        "withLabels": "WITH",
        "withoutLabels": "WITHOUT",
        "COUNT": "UNIT_COUNT",
        "DAYS": "UNIT_DAY",
        "none": "NONE",
        "nothing": "NONE",
        "Parameters nDaysSinceLastPull is too large": "DAYS_LARGE",
        "Parameters nDaysSinceLastPush is too large": "DAYS_LARGE",
        "Parameters latestPushedK is too large": "COUNT_LARGE",
        "Parameters latestPulledN is too large": "COUNT_LARGE"
    };

    constructor(
        private http: HttpClient,
    ) {
    }

    getI18nKey(str: string): string {
        if (this.I18nMap[str.trim()]) {
            return "TAG_RETENTION." + this.I18nMap[str.trim()];
        }
        return str;
    }

    getRetentionMetadata(): Observable<RuleMetadate> {
        return this.http.get(`/api/retentions/metadatas`)
            .pipe(map(response => response as RuleMetadate))
            .pipe(catchError(error => observableThrowError(error)));
    }

    getRetention(retentionId): Observable<Retention> {
        return this.http.get(`/api/retentions/${retentionId}`)
            .pipe(map(response => response as Retention))
            .pipe(catchError(error => observableThrowError(error)));
    }

    createRetention(retention: Retention) {
        return this.http.post(`/api/retentions`, retention)
            .pipe(catchError(error => observableThrowError(error)));
    }

    updateRetention(retentionId, retention: Retention) {
        return this.http.put(`/api/retentions/${retentionId}`, retention)
            .pipe(catchError(error => observableThrowError(error)));
    }

    getProjectInfo(projectId) {
        return this.http.get(`/api/projects/${projectId}`)
            .pipe(map(response => response as Project))
            .pipe(catchError(error => observableThrowError(error)));
    }

    runNowTrigger(retentionId) {
        return this.http.post(`/api/retentions/${retentionId}/executions`, {dry_run: false})
            .pipe(catchError(error => observableThrowError(error)));
    }

    whatIfRunTrigger(retentionId) {
        return this.http.post(`/api/retentions/${retentionId}/executions`, {dry_run: true})
            .pipe(catchError(error => observableThrowError(error)));
    }

    AbortRun(retentionId, executionId) {
        return this.http.patch(`/api/retentions/${retentionId}/executions/${executionId}`, {action: 'stop'})
            .pipe(catchError(error => observableThrowError(error)));
    }

    getRunNowList(retentionId, page: number, pageSize: number) {
        let params = new HttpParams();
        if (page && pageSize) {
            params = params.set('page', page + '').set('page_size', pageSize + '');
        }
        return this.http
          .get<HttpResponse<Array<any>>>(`/api/retentions/${retentionId}/executions`, buildHttpRequestOptionsWithObserveResponse(params))
          .pipe(catchError(error => observableThrowError(error)), );
    }

    getExecutionHistory(retentionId, executionId, page: number, pageSize: number) {
        let params = new HttpParams();
        if (page && pageSize) {
            params = params.set('page', page + '').set('page_size', pageSize + '');
        }
        return this.http.get<HttpResponse<Array<any>>>(`/api/retentions/${retentionId}/executions/${executionId}/tasks`,
            buildHttpRequestOptionsWithObserveResponse(params))
            .pipe(catchError(error => observableThrowError(error)));
    }

    seeLog(retentionId, executionId, taskId) {
        window.open(`api/retentions/${retentionId}/executions/${executionId}/tasks/${taskId}`, '_blank');
    }
}
