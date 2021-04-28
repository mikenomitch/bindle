import Component from '@glimmer/component';
import { action } from '@ember/object';
import fetch from 'fetch';

export default class Builder extends Component {
  vars = {};

  @action
  async generate() {
    const data = {
      variables: this.vars,
      name: this.args.packageName,
    };

    const url = 'http://localhost:9000/deploy';

    await fetch(url, {
      headers: {
        Accept: 'application/json',
        'Content-Type': 'application/json',
      },
      method: 'POST',
      body: JSON.stringify(data),
    });

    alert('WEBHOOK SENT - AKA JOB DEPLOYED!');

    window.location = '/';
  }

  @action
  setVariable(key, val) {
    this.vars[key] = val;
  }
}
