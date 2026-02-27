# frontend-design

This skill guides creation of distinctive, production-grade frontend interfaces
that avoid generic "AI slop" aesthetics. Implement real working code with
exceptional attention to aesthetic details and creative choices.

The user provides frontend requirements: a component, page, application, or
interface to build. They may include context about the purpose, audience, or
technical constraints.

## Design Thinking

Before coding, understand the context and commit to a BOLD aesthetic direction:

- Purpose: What problem does this interface solve? Who uses it?
- Tone: Pick an extreme: brutally minimal, maximalist chaos, retro-futuristic,
  organic/natural, luxury/refined, playful/toy-like, editorial/magazine, brutalist/raw,
  art deco/geometric, soft/pastel, industrial/utilitarian, etc.
- Constraints: Technical requirements (framework, performance, accessibility).
- Differentiation: What makes this UNFORGETTABLE?

CRITICAL: Choose a clear conceptual direction and execute it with precision. Bold
maximalism and refined minimalism both work — the key is intentionality, not intensity.

Then implement working code (HTML/CSS/JS, React, Vue, etc.) that is:

- Production-grade and functional
- Visually striking and memorable
- Cohesive with a clear aesthetic point-of-view
- Meticulously refined in every detail

## Frontend Aesthetics Guidelines

Focus on:

- **Typography**: Choose fonts that are beautiful, unique, and interesting. Avoid generic
  fonts like Arial and Inter; opt for distinctive choices that elevate aesthetics.
  Pair a distinctive display font with a refined body font.
- **Color & Theme**: Commit to a cohesive aesthetic. Use CSS variables for consistency.
  Dominant colors with sharp accents outperform timid, evenly-distributed palettes.
- **Motion**: Use animations for effects and micro-interactions. Focus on high-impact
  moments: one well-orchestrated page load with staggered reveals creates more delight
  than scattered micro-interactions.
- **Spatial Composition**: Unexpected layouts. Asymmetry. Overlap. Diagonal flow.
  Grid-breaking elements. Generous negative space OR controlled density.
- **Backgrounds & Visual Details**: Create atmosphere and depth rather than defaulting
  to solid colors. Use gradient meshes, noise textures, geometric patterns, layered
  transparencies, dramatic shadows, decorative borders.

NEVER use generic AI-generated aesthetics like overused font families (Inter, Roboto,
Arial, system fonts), cliched color schemes (particularly purple gradients on white
backgrounds), predictable layouts, and cookie-cutter design that lacks context-specific
character.

## Constraints

- NEVER use emojis anywhere in the UI — use icon libraries instead (react-icons)
- ALWAYS source icons from react-icons (FiX from 'react-icons/fi', MdX from 'react-icons/md',
  RiX from 'react-icons/ri', etc.) — never inline SVGs unless absolutely unavoidable
- Use CSS variables / Tailwind tokens for all colors, never hardcoded hex values
- Prioritize accessibility: proper ARIA labels, keyboard navigation, color contrast
- Every interaction must have a visual feedback state (hover, active, disabled, loading)
