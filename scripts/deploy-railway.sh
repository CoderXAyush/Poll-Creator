#!/bin/bash

echo "🚀 Setting up Railway deployment for Poll Creator..."

# Check if Railway CLI is installed
if ! command -v railway &> /dev/null; then
    echo "📦 Installing Railway CLI..."
    npm install -g @railway/cli
fi

# Login to Railway
echo "🔐 Please login to Railway..."
railway login

# Initialize Railway project
echo "🎯 Initializing Railway project..."
railway init

# Set environment variables
echo "⚙️ Setting environment variables..."
railway variables set PORT=8080
railway variables set NODE_ENV=production

# Deploy the application
echo "🚀 Deploying your application..."
railway up

# Get the deployment URL
echo "✅ Getting deployment URL..."
URL=$(railway status --json | jq -r '.deployments[0].url')

echo ""
echo "🎉 Deployment complete!"
echo "🌐 Your Poll Creator app is live at: $URL"
echo ""
echo "📝 Next steps:"
echo "   1. Visit your app URL to test it"
echo "   2. Share the URL with others to collect votes"
echo "   3. Monitor your app in the Railway dashboard"
echo ""