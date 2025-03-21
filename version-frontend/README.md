# Version Frontend

A modern React application built with Next.js that displays system information collected by the osquery backend service.

## Features

- Real-time display of system information:
  - OS Version details
  - Osquery Version
  - List of installed applications
- Modern, responsive UI built with Tailwind CSS
- Error handling and loading states
- Automatic data refresh
- Cross-browser compatibility

## Prerequisites

- Node.js 18.x or later
- npm or yarn
- Backend service running on port 7070

## Quick Start

1. Install dependencies:
```bash
npm install
# or
yarn install
```

2. Start the development server:
```bash
npm run dev
# or
yarn dev
```

3. Open [http://localhost:3000](http://localhost:3000) in your browser

## Development

### Project Structure

```
version-frontend/
├── app/                  # Next.js app directory
│   ├── components/      # React components
│   ├── lib/            # Utility functions
│   └── page.tsx        # Main page component
├── public/             # Static assets
└── styles/            # Global styles
```

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint
- `npm run test` - Run tests

### Environment Variables

Create a `.env.local` file in the root directory:

```env
NEXT_PUBLIC_API_URL=http://localhost:7070
```

## API Integration

The frontend communicates with the backend service through the following endpoint:

- `GET /api/latest_data` - Fetches the latest system information

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License
