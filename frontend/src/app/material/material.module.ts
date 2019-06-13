import {NgModule} from '@angular/core';
import {MatIconModule} from '@angular/material/icon'
import {MatButtonModule, MatDividerModule, MatFormFieldModule, MatInputModule} from '@angular/material'

@NgModule({
  declarations: [],
  imports: [
    MatInputModule,
    MatFormFieldModule,
    MatIconModule,
    MatButtonModule,
    MatDividerModule,
  ],
  exports: [
    MatInputModule,
    MatFormFieldModule,
    MatIconModule,
    MatButtonModule,
    MatDividerModule,
  ]
})
export class MaterialModule {}
