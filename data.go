package main

// Localization database
var localizationDB = map[string]map[string]string{
	"en": {
		"welcome_title":      "Welcome to Our App",
		"welcome_subtitle":   "Your journey starts here",
		"login_button":       "Log In",
		"signup_button":      "Sign Up",
		"navigation_home":    "Home",
		"navigation_about":   "About",
		"navigation_contact": "Contact",
		"footer_copyright":   "© 2024 Our Company. All rights reserved.",
		"user_profile_title": "User Profile",
		"user_profile_edit":  "Edit Profile",
		"settings_title":     "Settings",
		"settings_language":  "Language",
		"settings_theme":     "Theme",
		"error_404":          "Page not found",
		"error_500":          "Internal server error",
	},
	"es": {
		"welcome_title":      "Bienvenido a Nuestra App",
		"welcome_subtitle":   "Tu viaje comienza aquí",
		"login_button":       "Iniciar Sesión",
		"signup_button":      "Registrarse",
		"navigation_home":    "Inicio",
		"navigation_about":   "Acerca de",
		"navigation_contact": "Contacto",
		"footer_copyright":   "© 2024 Nuestra Empresa. Todos los derechos reservados.",
		"user_profile_title": "Perfil de Usuario",
		"user_profile_edit":  "Editar Perfil",
		"settings_title":     "Configuración",
		"settings_language":  "Idioma",
		"settings_theme":     "Tema",
		"error_404":          "Página no encontrada",
		"error_500":          "Error interno del servidor",
	},
	"fr": {
		"welcome_title":      "Bienvenue dans Notre App",
		"welcome_subtitle":   "Votre voyage commence ici",
		"login_button":       "Se Connecter",
		"signup_button":      "S'inscrire",
		"navigation_home":    "Accueil",
		"navigation_about":   "À Propos",
		"navigation_contact": "Contact",
		"footer_copyright":   "© 2024 Notre Entreprise. Tous droits réservés.",
		"user_profile_title": "Profil Utilisateur",
		"user_profile_edit":  "Modifier le Profil",
		"settings_title":     "Paramètres",
		"settings_language":  "Langue",
		"settings_theme":     "Thème",
		"error_404":          "Page non trouvée",
		"error_500":          "Erreur interne du serveur",
	},
	"de": {
		"welcome_title":      "Willkommen in Unserer App",
		"welcome_subtitle":   "Ihre Reise beginnt hier",
		"login_button":       "Anmelden",
		"signup_button":      "Registrieren",
		"navigation_home":    "Startseite",
		"navigation_about":   "Über Uns",
		"navigation_contact": "Kontakt",
		"footer_copyright":   "© 2024 Unser Unternehmen. Alle Rechte vorbehalten.",
		"user_profile_title": "Benutzerprofil",
		"user_profile_edit":  "Profil Bearbeiten",
		"settings_title":     "Einstellungen",
		"settings_language":  "Sprache",
		"settings_theme":     "Design",
		"error_404":          "Seite nicht gefunden",
		"error_500":          "Interner Serverfehler",
	},
}

// Component templates
var componentTemplates = map[string]ComponentTemplate{
	"welcome": {
		ComponentName: "WelcomeComponent",
		ComponentType: "functional",
		Template: `
import React from 'react';

const WelcomeComponent = ({ className = "welcome-container" }) => {
  return (
    <div className={className}>
      <div className="welcome-wrapper">
        <header className="welcome-header">
          <h1 className="welcome-title" data-l10n="welcome_title">
            {l10n.welcome_title}
          </h1>
          <p className="welcome-subtitle" data-l10n="welcome_subtitle">
            {l10n.welcome_subtitle}
          </p>
        </header>
        <div className="welcome-actions">
          <button 
            className="btn btn-primary" 
            onClick={() => {}}
            data-l10n="login_button"
          >
            {l10n.login_button}
          </button>
          <button 
            className="btn btn-secondary" 
            onClick={() => {}}
            data-l10n="signup_button"
          >
            {l10n.signup_button}
          </button>
        </div>
      </div>
    </div>
  );
};

export default WelcomeComponent;
`,
		RequiredKeys: []string{"welcome_title", "welcome_subtitle", "login_button", "signup_button"},
	},
	"navigation": {
		ComponentName: "NavigationComponent",
		ComponentType: "functional",
		Template: `
import React from 'react';

const NavigationComponent = ({ className = "navigation-container" }) => {
  return (
    <nav className={className}>
      <ul className="nav-list">
        <li className="nav-item">
          <a href="/" className="nav-link" data-l10n="navigation_home">
            {l10n.navigation_home}
          </a>
        </li>
        <li className="nav-item">
          <a href="/about" className="nav-link" data-l10n="navigation_about">
            {l10n.navigation_about}
          </a>
        </li>
        <li className="nav-item">
          <a href="/contact" className="nav-link" data-l10n="navigation_contact">
            {l10n.navigation_contact}
          </a>
        </li>
      </ul>
    </nav>
  );
};

export default NavigationComponent;
`,
		RequiredKeys: []string{"navigation_home", "navigation_about", "navigation_contact"},
	},
	"user_profile": {
		ComponentName: "UserProfileComponent",
		ComponentType: "functional",
		Template: `
import React from 'react';

const UserProfileComponent = ({ className = "user-profile-container" }) => {
  return (
    <div className={className}>
      <div className="profile-wrapper">
        <h2 className="profile-title" data-l10n="user_profile_title">
          {l10n.user_profile_title}
        </h2>
        <div className="profile-actions">
          <button 
            className="btn btn-outline" 
            onClick={() => {}}
            data-l10n="user_profile_edit"
          >
            {l10n.user_profile_edit}
          </button>
        </div>
      </div>
    </div>
  );
};

export default UserProfileComponent;
`,
		RequiredKeys: []string{"user_profile_title", "user_profile_edit"},
	},
	"footer": {
		ComponentName: "FooterComponent",
		ComponentType: "functional",
		Template: `
import React from 'react';

const FooterComponent = ({ className = "footer-container" }) => {
  return (
    <footer className={className}>
      <div className="footer-content">
        <p className="footer-copyright" data-l10n="footer_copyright">
          {l10n.footer_copyright}
        </p>
      </div>
    </footer>
  );
};

export default FooterComponent;
`,
		RequiredKeys: []string{"footer_copyright"},
	},
}

