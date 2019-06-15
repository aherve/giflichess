import { Injectable } from '@angular/core';

export interface AnalyticsEvent {
  eventCategory: string
  eventAction: string
  eventLabel?: string
  eventValue?: number
}

@Injectable({
  providedIn: 'root'
})
export class AnalyticsService {

  constructor() { }

  public trackEvent(evt: AnalyticsEvent) {
    if ("ga" in window && typeof window["ga"] === "function") {
      const ga = window["ga"] as Function
      ga('send', {
        hitType: 'event',
        ...evt,
      })
    }
  }
}
