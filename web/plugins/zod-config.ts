import { titleCase } from 'scule'
import { z } from 'zod'

export default defineNuxtPlugin(() => {
  z.config({
    customError: (iss) => {
      console.log(iss)
      if (iss.input === undefined && iss.code === 'invalid_type') {
        return `${titleCase(iss.path?.join('.') ?? 'This field')} is required`
      }
    },
  })
})
