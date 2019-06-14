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
<!-- Global site tag (gtag.js) - Google Analytics -->
<script async src="https://www.googletagmanager.com/gtag/js?id=${environment.analyticsToken}"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());

  gtag('config', '${environment.analyticsToken}');
</script>

    `
    document.body.appendChild(script)
  }
}
