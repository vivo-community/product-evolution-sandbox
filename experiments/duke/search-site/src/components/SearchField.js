import React, { Component } from 'react'

export default class SearchField extends Component {

  constructor(props) {
    super(props)
    this.shouldSetInputTextToDefaultValue = this.shouldSetInputTextToDefaultValue.bind(this)
  }

  // NOTE: jeeesh, this was convoluted and took a while to figure out (see below)
  // https://discuss.reactjs.org/t/how-to-pass-in-initial-value-to-form-fields/869/5
  // I just needed fields that I could initialize from query params - but also
  // be able to edit
  //
  shouldSetInputTextToDefaultValue (props) {
    let result = (this.previousDefaultValue != props.defaultValue) || (this.previousChangeIndicator != props.changeIndicator)
    return result
  }

  componentWillUpdate (nextProps, nextState) {
    let defaultText = nextProps.defaultValue
    let changeIndicator = nextProps.changeIndicator
      
    if (this.shouldSetInputTextToDefaultValue(nextProps))  {
      // set the default text input value if either the defaultText or the changeIndicator change 
      // set the default value
      let theInput = this.refs.myInput
      theInput.value = defaultText

      // save the default value and change indicator for later comparison
      this.previousDefaultText = defaultText
      this.previousChangeIndicator = changeIndicator
    }
  }

  shouldComponentUpdate (nextProps, nextState) {
    return this.shouldSetInputTextToDefaultValue(nextProps)
  }


  focus() {
    this.refs.myInput.focus()
  }

  // NOTE: ends up being a wrapper for typical input.value access
  // for a text input field - but in this case it's a <SearchField> component
  //
  set value(val) {
    const input = this.refs.myInput
    input.value = val
  }

  get value() {
    const input = this.refs.myInput
    const inputValue = input.value
    return inputValue
  }

  handleBlur(e) {
    // FIXME: sometimes getting error "cannot read property 'setState' of null"
    // ergo if(this)
    if(this) {
      this.setState({value: e.target.value})
    }
  }

  render() {
    const { label, placeholder, defaultValue, autofocus } = this.props

    let input = (
      <input onBlur = {this.handleBlur} type="text" ref="myInput" className="form-control input-sm"  placeholder={placeholder}/>
    )
    
    if (autofocus) {
      input = (
        <input autoFocus={true} onBlur = {this.handleBlur} type="text" ref="myInput" className="form-control input-sm"  placeholder={placeholder}/>
      )
    }

    return (
         <div className="form-group">
            <label>{label}</label>
              {input}
         </div>
    )

  }

}
