import { Component } from '@angular/core';

@Component({
  selector: 'ah-footer',
  templateUrl: './footer.component.html',
  styleUrls: ['./footer.component.scss']
})
export class FooterComponent {

  constructor() { }

  public get year () {
    return new Date().getFullYear()
  }

}
