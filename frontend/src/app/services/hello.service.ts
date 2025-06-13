import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { backendHost } from '../app.config';

@Injectable({
  providedIn: 'root'
})
export class HelloService {
  private apiUrl = `http://${backendHost}/api/v1/hello`;

  constructor(private http: HttpClient) {}

  getHello(): Observable<{ message: string }> {
    return this.http.get<{ message: string }>(this.apiUrl);
  }
}
