import { Component, OnInit } from '@angular/core';
import { Meta } from '@angular/platform-browser'
import {FormGroup, Validators, FormBuilder} from '@angular/forms';
import {AnalyticsService} from 'src/app/analytics.service';

@Component({
  selector: 'ah-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  public content = "Turn you lichess games into animated gifs"
  public isLoading = false
  public form: FormGroup

  constructor(
    private meta: Meta,
    private fb: FormBuilder,
    private analytics: AnalyticsService,
  ) { }

  public ngOnInit () {
    this.meta.addTags([
      {name: 'description', content: this.content}
    ])

    this.form = this.fb.group({
      lichessID: ['', [ Validators.required, Validators.minLength(8) ]]
    })
  }

  public submit() {
    if (!this.form.valid) { return }
    this.isLoading = true
    this.analytics.trackEvent({
      eventCategory: 'lichess',
      eventAction: 'loadGIF',
      eventLabel: this.form.value.lichessID,
    })
    const strip = this.form.value.lichessID
      .replace('http://', '')
      .replace('https://', '')
      .replace('lichess.org', '')
      .replace('/', '')
    window.location.href=`/lichess/${strip}`
  }

}
