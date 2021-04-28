import Route from '@ember/routing/route';
import fetch from 'fetch';
import { action } from '@ember/object';

export default class PackageRoute extends Route {
  async model(params) {
    const response = await fetch(
      `http://localhost:9000/package?name=${params.package_name}`
    );
    const p = await response.json();
    p.name = params.package_name;
    return { package: p };
  }
}
