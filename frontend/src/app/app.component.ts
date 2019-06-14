import { Component } from '@angular/core';
import {AnalyticsService} from 'src/app/analytics.service';

@Component({
  selector: 'ah-root',
  template: `<ah-layout></ah-layout>`,
  styles: []
})
export class AppComponent {
  constructor(_: AnalyticsService) {}
}
