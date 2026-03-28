interface InputProps {
  label: string
  name: string
  type?: string
  required?: boolean
  placeholder?: string
  defaultValue?: string
}

export function InputField({ label, name, type = 'text', required = true, placeholder, defaultValue }: InputProps) {
  return (
    <label className="field">
      <span>{label}</span>
      <input name={name} type={type} required={required} placeholder={placeholder} defaultValue={defaultValue} />
    </label>
  )
}

interface SelectProps {
  label: string
  name: string
  options: Array<{ value: string; label: string }>
}

export function SelectField({ label, name, options }: SelectProps) {
  return (
    <label className="field">
      <span>{label}</span>
      <select name={name}>
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    </label>
  )
}

export function TextareaField({ label, name, required = false }: Omit<InputProps, 'type'>) {
  return (
    <label className="field">
      <span>{label}</span>
      <textarea name={name} required={required} rows={3} />
    </label>
  )
}
