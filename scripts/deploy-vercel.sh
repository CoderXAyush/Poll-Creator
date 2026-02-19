#!/bin/bash

echo "🚀 Setting up Vercel deployment for Poll Creator..."

# Check if Vercel CLI is installed
if ! command -v vercel &> /dev/null; then
    echo "📦 Installing Vercel CLI..."
    npm install -g vercel
fi

# Login to Vercel
echo "🔐 Please login to Vercel..."
vercel login

# Initialize Vercel project
echo "🎯 Setting up Vercel project..."
vercel --confirm

# Deploy the application
echo "🚀 Deploying your application..."
vercel --prod

echo ""
echo "🎉 Deployment complete!"
echo ""
echo "📝 Your Poll Creator app is now live on Vercel!"
echo "🌐 You can view it at the URL provided above"
echo ""
echo "📋 Next steps:"
echo "   1. Visit your app URL to test it"
echo "   2. Share the URL with others to collect votes"
echo "   3. Monitor your app in the Vercel dashboard"
echo "   4. Set up automatic deployments in Vercel dashboard by connecting your GitHub repo"
echo ""