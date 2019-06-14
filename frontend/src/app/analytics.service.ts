import { Injectable } from '@angular/core';
import { environment } from '../environments/environment'

@Injectable({
  providedIn: 'root'
})
export class AnalyticsService {

  constructor() {
    this.loadGaScript()
  }

  private loadGaScript () {
    if (!('analyticsToken' in environment)) { return }
    const script = document.createElement('script')
    script.text = `
      (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
      (i[r].q=i[r].q||[]).push(arguments)
      },i[r].l=1*new Date();a=s.createElement(o),
      m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
      })(window,document,'script','https://www.google-analytics.com/analytics.js','ga');
      ga('create', '${environment["analyticsToken"]}', 'auto');
    `
    document.body.appendChild(script)
  }
}
