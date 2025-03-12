import { Formik, Form, Field, FormikHelpers } from "formik"
import { toFormikValidationSchema } from "zod-formik-adapter"
import { registerSchema, type RegisterInput } from "../lib/validation"
import { Link, useNavigate } from "react-router-dom"
import { useState } from "react"
import toast from "react-hot-toast"

export default function Register() {

  }

  return (
    <div className="min-h-screen bg-base-200 flex items-center justify-center p-4">
      <div className="card w-full max-w-md bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title text-2xl font-bold text-center mb-2">
            Create Account
          </h2>
          <p className="text-center text-base-content/70 mb-6">
            Join us today! Fill in your details below
          </p>

          <Formik
            initialValues={{
              username: "",
              email: "",
              password: "",
            }}
            validationSchema={toFormikValidationSchema(registerSchema)}
            onSubmit={handleSubmit}
          >
            {({ errors, touched, isSubmitting }) => (
              <Form className="space-y-4">
                <div className="form-control w-full">
                  <label className="label">
                    <span className="label-text">Username</span>
                  </label>
                  <Field
                    name="username"
                    type="text"
                    placeholder="johndoe"
                    className="input input-bordered w-full"
                  />
                  {errors.username && touched.username && (
                    <label className="label">
                      <span className="label-text-alt text-error">
                        {errors.username}
                      </span>
                    </label>
                  )}
                </div>

                <div className="form-control w-full">
                  <label className="label">
                    <span className="label-text">Email</span>
                  </label>
                  <Field
                    name="email"
                    type="email"
                    placeholder="your@email.com"
                    className="input input-bordered w-full"
                  />
                  {errors.email && touched.email && (
                    <label className="label">
                      <span className="label-text-alt text-error">
                        {errors.email}
                      </span>
                    </label>
                  )}
                </div>

                <div className="form-control w-full">
                  <label className="label">
                    <span className="label-text">Password</span>
                  </label>
                  <Field
                    name="password"
                    type="password"
                    placeholder="••••••••"
                    className="input input-bordered w-full"
                  />
                  {errors.password && touched.password && (
                    <label className="label">
                      <span className="label-text-alt text-error">
                        {errors.password}
                      </span>
                    </label>
                  )}
                </div>
                {backendError && (
                  <div className="alert alert-error">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      className="stroke-current shrink-0 h-6 w-6"
                      fill="none"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                    <span>{backendError}</span>
                  </div>
                )}
                <button
                  type="submit"
                  className={`btn btn-primary w-full ${
                    isSubmitting ? "loading" : ""
                  }`}
                  disabled={isSubmitting}
                >
                  {isSubmitting ? "Creating Account..." : "Create Account"}
                </button>
              </Form>
            )}
          </Formik>

          <div className="divider my-6">OR</div>

          <div className="flex flex-col gap-3">
            <button className="btn btn-outline">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                className="h-5 w-5 mr-2"
              >
                <path
                  d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                  fill="#4285F4"
                />
                <path
                  d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                  fill="#34A853"
                />
                <path
                  d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                  fill="#FBBC05"
                />
                <path
                  d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                  fill="#EA4335"
                />
              </svg>
              Sign up with Google
            </button>
          </div>

          <div className="text-center mt-6">
            <p>
              Already have an account?
              <Link to="/login" className="link link-primary ml-1">
                Sign in
              </Link>
            </p>
          </div>

          <div className="text-center text-sm mt-4 text-base-content/70">
            <p>By signing up, you agree to our</p>
            <div className="flex justify-center gap-2 mt-1">
              <Link to="/terms" className="link link-hover">
                Terms of Service
              </Link>
              <span>&</span>
              <Link to="/privacy" className="link link-hover">
                Privacy Policy
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
