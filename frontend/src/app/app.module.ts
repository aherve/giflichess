import {BrowserModule} from '@angular/platform-browser';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {NgModule} from '@angular/core';

import {AppRoutingModule} from './app-routing.module';
import {AppComponent} from './app.component';
import {HomeComponent} from './components/home/home.component';
import {FlexLayoutModule} from '@angular/flex-layout';
import {FooterComponent} from './components/footer/footer.component';
import {LayoutComponent} from './components/layout/layout.component'
import {MaterialModule} from 'src/app/material/material.module';
import {FormsModule, ReactiveFormsModule} from '@angular/forms';
import {AnalyticsService} from 'src/app/analytics.service';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    FooterComponent,
    LayoutComponent
  ],
  imports: [
    ReactiveFormsModule,
    FormsModule,
    MaterialModule,
    BrowserAnimationsModule,
    FlexLayoutModule,
    BrowserModule,
    AppRoutingModule
  ],
  providers: [AnalyticsService],
  bootstrap: [AppComponent]
})
export class AppModule { }
