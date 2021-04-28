import Route from '@ember/routing/route';
import fetch from 'fetch';

export default class CatalogRoute extends Route {
  async model() {
    const response = await fetch('http://localhost:9000/catalog');
    return await response.json();
  }
}
