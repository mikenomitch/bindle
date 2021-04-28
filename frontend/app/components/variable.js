import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

export default class CounterComponent extends Component {
  @tracked inputVal = this.args.default;

  // onRender
  @action
  setVal(i) {
    this.args.setVariable(this.args.name, i.target.value);
  }

  @action
  setDefaultValue() {
    this.args.setVariable(this.args.name, this.args.default);
  }
}
