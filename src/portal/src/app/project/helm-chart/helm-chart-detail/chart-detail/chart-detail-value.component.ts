import {
  Component,
  Input,
  OnInit,
  ChangeDetectionStrategy
} from "@angular/core";

@Component({
  selector: "hbr-chart-detail-value",
  templateUrl: "./chart-detail-value.component.html",
  styleUrls: ["./chart-detail-value.component.scss"],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class ChartDetailValueComponent implements OnInit {
  @Input() values;
  @Input() yaml;

  // Default set to yaml file
  valueMode = false;
  valueHover = false;
  yamlHover = true;

  objKeys = Object.keys;

  constructor() {}

  ngOnInit(): void {
  }

  public get isValueMode() {
    return this.valueMode;
  }

  isHovering(view: string) {
    if (view === 'value') {
      return this.valueHover ? true : false;
    } else {
      return this.yamlHover ? true : false;
    }
  }

  showYamlFile(showYaml: boolean) {
    this.valueMode = !showYaml;
  }

  mouseEnter(mode: string) {
    if (mode === "value") {
      this.valueHover = true;
    } else {
      this.yamlHover = true;
    }
  }

  mouseLeave(mode: string) {
    if (mode === "value") {
      this.valueHover = false;
    } else {
      this.yamlHover = false;
    }
  }
}
