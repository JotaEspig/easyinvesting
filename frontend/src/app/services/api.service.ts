import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { apiUrl } from '../app.config';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  constructor(private http: HttpClient) {}

  getRequest<T>(endpoint: string, options?: object): Observable<T> {
    console.log(`GET request to: ${apiUrl + endpoint}`);
    return this.http.get<T>(apiUrl + endpoint, options);
  }

  postRequest<T>(endpoint: string, body: any, options?: object): Observable<T> {
    console.log(`POST request to: ${apiUrl + endpoint}`);
    return this.http.post<T>(apiUrl + endpoint, body, options);
  }
}
