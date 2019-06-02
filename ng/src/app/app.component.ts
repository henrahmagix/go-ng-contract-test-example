import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import * as getHelloContract from './contract/get-hello.json';
import { Response as HelloGetResponse } from './contract/get-hello.js';

@Component({
  selector: 'app-root',
  template: `
    <ng-container *ngIf="ready">
      Response = "{{response}}"
    </ng-container>
  `,
})
export class AppComponent implements OnInit {
  response;
  ready = false;

  constructor(
    private http: HttpClient,
  ) {
  }
  
  ngOnInit() {
    const url = `http://localhost:8080${getHelloContract.path}`;
    this.http.request<HelloGetResponse>(getHelloContract.method, url).forEach(res => {
      this.response = res.msg;
      console.log('res', res);
      this.ready = true;
    });
  }
}
