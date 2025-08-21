import { Injectable } from '@angular/core';
import { AngularFireStorage } from '@angular/fire/compat/storage';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class FirebaseService {
  imageURL = new BehaviorSubject<string | null>(null);
  imageURL$ = this.imageURL.asObservable();
  constructor(private storage: AngularFireStorage) {}

  uploadImage = (selectedImages: File[]): Promise<string> => {
    const file = selectedImages[0];
    const filePath = `postIMG/${Date.now()}_${file.name}`;
    const fileRef = this.storage.ref(filePath);
    const metadata = { contentType: file.type };

    return this.storage.upload(filePath, file, metadata).then(() => {
      return fileRef.getDownloadURL().toPromise(); // returns a single string (URL)
    });
  };
}
