import { Router } from '@angular/router';
import { Component, OnInit, Input } from '@angular/core';

import { HelmChartSearchResultItem, HelmChartVersion, HelmChartMaintainer } from '../../project/helm-chart/helm-chart.interface.service';

import { SearchTriggerService } from '../../base/global-search/search-trigger.service';
import { ProjectService } from '@harbor/ui';

@Component({
  selector: 'list-chart-version-ro',
  templateUrl: './list-chart-version-ro.component.html'
})
export class ListChartVersionRoComponent implements OnInit {

  @Input() projectId: number;
  @Input() charts: HelmChartSearchResultItem;

  constructor(
    private searchTrigger: SearchTriggerService,
    private projectService: ProjectService,
    private router: Router) { }

  ngOnInit() {}

  getStatusString(chart: HelmChartVersion) {
    if (chart.deprecated) {
      return "HELM_CHART.DEPRECATED";
    } else {
      return "HELM_CHART.ACTIVE";
    }
  }

  getMaintainerString(maintainers: HelmChartMaintainer[]) {
    if (!maintainers || maintainers.length < 1) {
      return "";
    }

    let maintainer_string = maintainers[0].name;
    if (maintainers.length > 1) {
      maintainer_string = "HELM_CHART.OTHER_MAINTAINERS";
    }
    return maintainer_string;
  }

  getMaintainerTranslateInfo(maintainers: HelmChartMaintainer[]) {
    if (!maintainers || maintainers.length < 1) {
      return {};
    }
    let name = maintainers[0].name;
    let number = maintainers.length;
    return {name: name, number: number};
  }

  gotoChartVersion(chartVersion: HelmChartVersion) {
    this.searchTrigger.closeSearch(true);
    let [projectName, chartName] = chartVersion.name.split('/');
    this.projectService.listProjects(projectName).subscribe( res => {
      let projects = res.body || [];
      if (projects || projects.length >= 1) {
        let linkUrl = ['harbor', 'projects', projects[0].project_id, 'helm-charts', chartName, 'versions', chartVersion.version];
        this.router.navigate(linkUrl);
      } else {
        return;
      }
    });
  }
}
