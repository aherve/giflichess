import { Component, OnInit } from '@angular/core';
import { Meta } from '@angular/platform-browser'
import {FormGroup, Validators, FormBuilder} from '@angular/forms';

@Component({
  selector: 'ah-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  public content = "Turn you lichess games into animated gifs"
  public isLoading = true
  public form: FormGroup

  constructor(
    private meta: Meta,
    private fb: FormBuilder,
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
    const strip = this.form.value.lichessID
      .replace('http://', '')
      .replace('https://', '')
      .replace('lichess.org', '')
      .replace('/', '')
    window.location.href=`/lichess/${strip}`
  }

}
